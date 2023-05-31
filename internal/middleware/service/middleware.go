package service

import (
	"context"
	"errors"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"tini-paas/internal/middleware/model"
	"tini-paas/internal/middleware/proto/middleware"
	"tini-paas/internal/middleware/repository"
	"tini-paas/pkg/common"
)

// MiddlewareService 中间件服务接口
type MiddlewareService interface {
	AddMiddleware(*model.Middleware) (int64, error)
	DeleteMiddleware(int64) error
	UpdateMiddleware(*model.Middleware) error
	FindMiddlewareByID(int64) (*model.Middleware, error)
	FindAllMiddleware() ([]model.Middleware, error)

	// FindAllMiddlewareByTypeID 根据类型查找中间件
	FindAllMiddlewareByTypeID(int64) ([]model.Middleware, error)

	CreateToK8s(*middleware.MiddlewareInfo) error
	DeleteFromK8s(*model.Middleware) error
	UpdateToK8s(*middleware.MiddlewareInfo) error
}

// NewMiddlewareService 初始化中间件服务
func NewMiddlewareService(middlewareRepository repository.MiddlewareRepository, clientSet *kubernetes.Clientset) MiddlewareService {
	return &MiddlewareDataService{
		MiddlewareRepository: middlewareRepository,
		K8sClientSet:         clientSet,
	}
}

// MiddlewareDataService  中间件服务操作对象
type MiddlewareDataService struct {
	MiddlewareRepository repository.MiddlewareRepository
	K8sClientSet         *kubernetes.Clientset
}

func (m *MiddlewareDataService) AddMiddleware(middle *model.Middleware) (int64, error) {
	return m.MiddlewareRepository.CreateMiddleware(middle)
}

func (m *MiddlewareDataService) DeleteMiddleware(i int64) error {
	return m.MiddlewareRepository.DeleteMiddleware(i)
}

func (m *MiddlewareDataService) UpdateMiddleware(middle *model.Middleware) error {
	return m.MiddlewareRepository.UpdateMiddleware(middle)
}

func (m *MiddlewareDataService) FindMiddlewareByID(i int64) (*model.Middleware, error) {
	return m.MiddlewareRepository.FindMiddlewareByID(i)
}

func (m *MiddlewareDataService) FindAllMiddleware() ([]model.Middleware, error) {
	return m.MiddlewareRepository.FindAll()
}

func (m *MiddlewareDataService) FindAllMiddlewareByTypeID(i int64) ([]model.Middleware, error) {
	return m.MiddlewareRepository.FindAllByTypeID(i)
}

func (m *MiddlewareDataService) CreateToK8s(info *middleware.MiddlewareInfo) error {
	// 设置stateful信息
	statefulSet := m.setStatefulSet(info)

	// 先查询是否存在
	_, err := m.K8sClientSet.AppsV1().StatefulSets(info.MiddleNamespace).Get(context.TODO(), info.MiddleName, v12.GetOptions{})
	if err != nil {
		// 之前不存在 -> 创建
		_, err = m.K8sClientSet.AppsV1().StatefulSets(info.MiddleNamespace).Create(context.TODO(), statefulSet, v12.CreateOptions{})
		if err != nil {
			// 创建失败
			common.Error(err)
			return err
		}

		// 创建成功
		common.Info("中间件：" + info.MiddleName + "创建成功")
		return nil
	}

	// 之前存在
	common.Error("中间件：" + info.MiddleName + "已经存在")
	return errors.New("中间件：" + info.MiddleName + "已经存在")
}

func (m *MiddlewareDataService) DeleteFromK8s(middle *model.Middleware) error {
	// 先删除k8s
	err := m.K8sClientSet.AppsV1().StatefulSets(middle.MiddleNamespace).Delete(context.TODO(), middle.MiddleName, v12.DeleteOptions{})
	if err != nil {
		// 删除失败
		common.Error(err)
		return err
	}

	// 删除数据库信息
	err = m.MiddlewareRepository.DeleteMiddleware(middle.ID)
	if err != nil {
		common.Error(err)
		return err
	}
	common.Info("删除中间件：" + middle.MiddleName + "成功！")

	return nil
}

func (m *MiddlewareDataService) UpdateToK8s(info *middleware.MiddlewareInfo) error {
	// 设置stateful信息
	statefulSet := m.setStatefulSet(info)

	// 先查询是否存在
	_, err := m.K8sClientSet.AppsV1().StatefulSets(info.MiddleNamespace).Get(context.TODO(), info.MiddleName, v12.GetOptions{})
	if err != nil {
		common.Error(err)
		return errors.New("中间件 " + info.MiddleName + " 不存在请先创建")
	}

	// 已经存在，可以更新
	_, err = m.K8sClientSet.AppsV1().StatefulSets(info.MiddleNamespace).Update(context.TODO(), statefulSet, v12.UpdateOptions{})
	if err != nil {
		// 更新失败
		common.Error(err)
		return err
	}

	// 更新成功
	common.Info("中间件 " + info.MiddleName + " 更新成功！")
	return nil
}

// setStatefulSet 根据info信息设置值
func (m *MiddlewareDataService) setStatefulSet(info *middleware.MiddlewareInfo) *v1.StatefulSet {
	statefulSet := &v1.StatefulSet{}

	// 设置接口类型
	statefulSet.TypeMeta = v12.TypeMeta{
		Kind:       "StatefulSet",
		APIVersion: "v1",
	}

	// 设置详细信息
	statefulSet.ObjectMeta = v12.ObjectMeta{
		Name:      info.MiddleName,
		Namespace: info.MiddleNamespace,
		Labels: map[string]string{
			"app-name": info.MiddleName,
			"author":   "Paas",
		},
	}

	statefulSet.Name = info.MiddleName
	statefulSet.Spec = v1.StatefulSetSpec{
		Replicas: &info.MiddleReplicas,
		Selector: &v12.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": info.MiddleName,
			},
		},

		// 设置容器模板
		Template: v13.PodTemplateSpec{
			ObjectMeta: v12.ObjectMeta{
				Labels: map[string]string{
					"app-name": info.MiddleName,
				},
			},

			// 设置容器详情
			Spec: v13.PodSpec{
				Containers: []v13.Container{
					{
						Name:  info.MiddleName,
						Image: info.MiddleDockerImageVersion,
						// 获取容器的端口
						Ports: m.getContainerPort(info),
						// 设置环境变量
						Env: m.getEnv(info),
						// 设置容器资源配额
						Resources: m.getResources(info),
						// 设置挂载目录
						VolumeMounts: m.setMounts(info),
					},
				},

				// 不能设置为0，不安全
				TerminationGracePeriodSeconds: m.getTime("10"),
				// 设置私有仓库密钥
				ImagePullSecrets: nil,
			},
		},

		// 获取PVC
		VolumeClaimTemplates: m.getPVC(info),
		ServiceName:          info.MiddleName,
	}

	return statefulSet
}

// getContainerPort 获取容器端口
func (m *MiddlewareDataService) getContainerPort(info *middleware.MiddlewareInfo) []v13.ContainerPort {
	var containerPort []v13.ContainerPort

	for _, port := range info.MiddlePort {
		containerPort = append(containerPort, v13.ContainerPort{
			Name:          "middle-port-" + strconv.FormatInt(int64(port.MiddlePort), 10),
			ContainerPort: port.MiddlePort,
			Protocol:      m.getProtocol(port.MiddleProtocol),
		})
	}
	return containerPort
}

// getProtocol 获取protocol协议
func (m *MiddlewareDataService) getProtocol(protocol string) v13.Protocol {
	switch protocol {
	case "TCP":
		return "TCP"
	case "UDP":
		return "UDP"
	case "SCTP":
		return "SCTP"
	default:
		return "TCP"
	}
}

// getEnv 获取中间件的环境变量
func (m *MiddlewareDataService) getEnv(info *middleware.MiddlewareInfo) []v13.EnvVar {
	var envVar []v13.EnvVar

	for _, env := range info.MiddleEnv {
		envVar = append(envVar, v13.EnvVar{
			Name:      env.EbvKey,
			Value:     env.EnvValue,
			ValueFrom: nil,
		})
	}
	return envVar
}

// getResources 获取容器的资源配额
func (m *MiddlewareDataService) getResources(info *middleware.MiddlewareInfo) v13.ResourceRequirements {
	source := v13.ResourceRequirements{}

	// 最大能够使用的资源
	source.Limits = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(info.MiddleCpu), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(info.MiddleMemory), 'f', 6, 64)),
	}
	// 最小请求资源
	source.Requests = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(info.MiddleCpu/2), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(info.MiddleMemory/2), 'f', 6, 64)),
	}
	return source
}

// getPVC 获取PVC
func (m *MiddlewareDataService) getPVC(info *middleware.MiddlewareInfo) []v13.PersistentVolumeClaim {
	var pvcAll []v13.PersistentVolumeClaim
	if len(info.MiddleStorage) == 0 {
		return pvcAll
	}

	for _, storage := range info.MiddleStorage {
		pvc := &v13.PersistentVolumeClaim{
			TypeMeta: v12.TypeMeta{
				Kind:       "PersistentVolumeClaim",
				APIVersion: "v1",
			},
			ObjectMeta: v12.ObjectMeta{
				Name:      storage.MiddleStorageName,
				Namespace: info.MiddleNamespace,
				Annotations: map[string]string{
					"pv.kubernetes.io/bound-by-controller":          "yes",
					"volume.beta,kubernetes.io/storage-provisioner": "rbd.csi.ceph.com",
				},
			},
			Spec: v13.PersistentVolumeClaimSpec{
				AccessModes:      m.getAccessModes(storage.MiddleStorageAccessMode),
				Resources:        m.getPVCResource(storage.MiddleStorageSize),
				VolumeName:       storage.MiddleStorageName,
				StorageClassName: &storage.MiddleStorageClass,
			},
		}
		pvcAll = append(pvcAll, *pvc)
	}
	return pvcAll
}

// getAccessModes 获取存储权限
func (m *MiddlewareDataService) getAccessModes(accessMode string) []v13.PersistentVolumeAccessMode {
	var pvm []v13.PersistentVolumeAccessMode

	var pm v13.PersistentVolumeAccessMode
	switch accessMode {
	case "ReadWriteOnce":
		pm = v13.ReadWriteOnce
	case "ReadOnlyMany":
		pm = v13.ReadOnlyMany
	case "ReadWriteMany":
		pm = v13.ReadWriteMany
	case "ReadWriteOncePod":
		pm = v13.ReadWriteOncePod
	default:
		pm = v13.ReadWriteOnce
	}
	pvm = append(pvm, pm)
	return pvm
}

// getPVCResource 获取PVC大小
func (m *MiddlewareDataService) getPVCResource(size float32) v13.ResourceRequirements {
	source := v13.ResourceRequirements{}

	source.Requests = v13.ResourceList{
		"storage": resource.MustParse(strconv.FormatFloat(float64(size), 'f', 6, 64) + "Gi"),
	}
	return source
}

// setMounts 设置存储路径
func (m *MiddlewareDataService) setMounts(info *middleware.MiddlewareInfo) []v13.VolumeMount {
	var mount []v13.VolumeMount
	if len(info.MiddleStorage) == 0 {
		// 没有挂载存储数据
		return mount
	}

	for _, storage := range info.MiddleStorage {
		mt := &v13.VolumeMount{
			Name:      storage.MiddleStorageName,
			MountPath: storage.MiddleStoragePath,
		}
		mount = append(mount, *mt)
	}
	return mount
}

// setTime 优雅终止时间
func (m *MiddlewareDataService) getTime(stringTime string) *int64 {
	parseInt, err := strconv.ParseInt(stringTime, 10, 64)
	if err != nil {
		common.Error(err)
		return nil
	}
	return &parseInt
}

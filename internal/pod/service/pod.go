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
	"tini-paas/internal/pod/model"
	"tini-paas/internal/pod/proto/pod"
	"tini-paas/internal/pod/repository"
	"tini-paas/pkg/common"
)

// PodService pod服务
type PodService interface {
	AddPod(*model.Pod) (int64, error)
	DeletedPod(int64) error
	UpdatePod(*model.Pod) error
	FindPodByID(int64) (*model.Pod, error)
	FindAllPod() ([]model.Pod, error)
	CreateToK8s(*pod.PodInfo) error
	UpdateToK8s(*pod.PodInfo) error
	DeletedFromK8s(*model.Pod) error
}

// PodDataService pod数据服务
type PodDataService struct {
	// PodRepository 操作数据库接口
	PodRepository repository.PodRepository

	// K8sClientSet k8s客户端集合
	K8sClientSet *kubernetes.Clientset

	// deployment 发布控制器
	deployment *v1.Deployment
}

// NewPodService 初始化pod服务
func NewPodService(podRepository repository.PodRepository, clientSet *kubernetes.Clientset) PodService {
	return &PodDataService{
		PodRepository: podRepository,
		K8sClientSet:  clientSet,
		deployment:    &v1.Deployment{},
	}
}

// AddPod 添加pod
func (p *PodDataService) AddPod(pod *model.Pod) (int64, error) {
	return p.PodRepository.CreatePod(pod)
}

// DeletedPod 删除pod
func (p *PodDataService) DeletedPod(podID int64) error {
	return p.PodRepository.DeletePodByID(podID)
}

// UpdatePod 更新pod
func (p *PodDataService) UpdatePod(pod *model.Pod) error {
	return p.PodRepository.UpdatePod(pod)
}

// FindPodByID 查找pod
func (p *PodDataService) FindPodByID(podID int64) (*model.Pod, error) {
	return p.PodRepository.FindPodByID(podID)
}

// FindAllPod 查找全部pod
func (p *PodDataService) FindAllPod() ([]model.Pod, error) {
	return p.PodRepository.FindAll()
}

// CreateToK8s 创建pod到k8s
func (p *PodDataService) CreateToK8s(info *pod.PodInfo) error {
	// 根据podInfo设置发布控制器Deployment
	p.SetDeployment(info)

	_, err := p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Get(context.TODO(), info.PodName, v12.GetOptions{})
	if err != nil {
		// 之前不存在此名称的pod -> 创建pod
		_, err = p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Create(context.TODO(), p.deployment, v12.CreateOptions{})
		if err != nil {
			// 创建失败
			//common1.Error(err)
			return err
		}
		common.Info("创建成功")
		return nil
	}

	// 没报错说明已经存在
	// TODO 业务逻辑
	common.Error("Pod " + info.PodName + "已经存在")
	return errors.New("Pod " + info.PodName + " 已经存在")
}

// UpdateToK8s 更新pod到k8s
func (p *PodDataService) UpdateToK8s(info *pod.PodInfo) error {
	// 根据podInfo设置发布控制器Deployment
	p.SetDeployment(info)

	_, err := p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Get(context.TODO(), info.PodName, v12.GetOptions{})
	if err != nil {
		// 之前不存在的pod -> 不更新
		common.Error(err)
		return errors.New("Pod " + info.PodName + " 不存在请先创建")
	}

	// 之前存在，可以更新
	_, err = p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Update(context.TODO(), p.deployment, v12.UpdateOptions{})
	if err != nil {
		// 更新失败
		//common.Error(err)
		return err
	}
	common.Info(info.PodName + " 更新成功")
	return nil
}

// DeletedFromK8s 从k8s删除pod
func (p *PodDataService) DeletedFromK8s(pod *model.Pod) error {
	err := p.K8sClientSet.AppsV1().Deployments(pod.PodNamespace).Delete(context.TODO(), pod.PodName, v12.DeleteOptions{})
	if err != nil {
		// 删除错误
		// TODO 业务逻辑
		//common.Error(err)
		return err
	}

	// 正常删除 -> 删除数据库数据
	err = p.DeletedPod(pod.ID)
	if err != nil {
		//common.Error(err)
		return err
	}
	common.Info("删除Pod ID: " + strconv.FormatInt(pod.ID, 10) + " 成功!")
	return nil
}

// SetDeployment 设置发布控制器
func (p *PodDataService) SetDeployment(info *pod.PodInfo) {
	deployment := &v1.Deployment{}

	// deployment元数据类型
	deployment.TypeMeta = v12.TypeMeta{
		Kind:       "deployment",
		APIVersion: "v1",
	}

	// deployment持久化目标元数据
	deployment.ObjectMeta = v12.ObjectMeta{
		Name:      info.PodName,
		Namespace: info.PodNamespace,
		Labels: map[string]string{
			"app-name": info.PodName,
		},
	}

	// pod名称
	deployment.Name = info.PodName

	// 详细应用参数
	deployment.Spec = v1.DeploymentSpec{
		//副本数量
		Replicas: &info.PodReplicas,

		// 通过标签进行匹配
		Selector: &v12.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": info.PodName,
			},
			MatchExpressions: nil,
		},

		// 容器模板
		Template: v13.PodTemplateSpec{
			// 目标元数据
			ObjectMeta: v12.ObjectMeta{
				Labels: map[string]string{
					"app-name": info.PodName,
				},
			},

			// pod详细信息
			Spec: v13.PodSpec{
				// 容器
				Containers: []v13.Container{
					{
						Name:            info.PodName,               // pod名称
						Image:           info.PodImage,              // pod镜像
						Ports:           p.getContainerPort(info),   // pod容器端口
						Env:             p.getEnv(info),             // pod环境变量
						Resources:       p.getResources(info),       // pod资源限制
						ImagePullPolicy: p.getImagePullPolicy(info), // pod镜像拉取策略
					},
				},
			},
		},
		Strategy:                v1.DeploymentStrategy{},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    nil,
		Paused:                  false,
		ProgressDeadlineSeconds: nil,
	}

	// 将配置信息赋值
	p.deployment = deployment
}

// getContainerPort 生成容器端口
func (p *PodDataService) getContainerPort(info *pod.PodInfo) []v13.ContainerPort {
	var containerPort []v13.ContainerPort
	for _, port := range info.PodPort {
		containerPort = append(containerPort, v13.ContainerPort{
			Name:          "port-" + strconv.FormatInt(int64(port.ContainerPort), 10),
			ContainerPort: port.ContainerPort,
			Protocol:      p.getProtocol(port.Protocol),
		})
	}
	return containerPort
}

// getProtocol 获取协议
func (p *PodDataService) getProtocol(protocol string) v13.Protocol {
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

// getEnv 环境变量
func (p *PodDataService) getEnv(info *pod.PodInfo) []v13.EnvVar {
	var envVar []v13.EnvVar

	for _, env := range info.PodEnv {
		envVar = append(envVar, v13.EnvVar{
			Name:      env.EnvKey,
			Value:     env.EnvValue,
			ValueFrom: nil,
		})
	}
	return envVar
}

// getResources 限制使用的最大资源
func (p *PodDataService) getResources(info *pod.PodInfo) v13.ResourceRequirements {
	var source v13.ResourceRequirements
	const base = 4.0

	// 最大能够使用的资源
	source.Limits = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(info.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(info.PodMemoryMax), 'f', 6, 64)),
	}

	// 满足最少使用的资源量
	source.Requests = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(info.PodCpuMax/base), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(info.PodMemoryMax/base), 'f', 6, 64)),
	}

	return source
}

// getImagePullPolicy 镜像拉取策略
func (p *PodDataService) getImagePullPolicy(info *pod.PodInfo) v13.PullPolicy {
	switch info.PodPullPolicy {
	case "Always":
		return "Always"
	case "Never":
		return "Never"
	case "IfNotPresent":
		return "IfNotPresent"
	default:
		return "Always"
	}
}

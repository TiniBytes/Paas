package service

import (
	"context"
	"errors"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"tini-paas/internal/svc/model"
	"tini-paas/internal/svc/proto/svc"
	"tini-paas/internal/svc/repository"
	"tini-paas/pkg/common"
)

// SvcService 服务
type SvcService interface {
	// AddSvc 添加service
	AddSvc(*model.Svc) (int64, error)

	// DeleteSvc 删除service
	DeleteSvc(int64) error

	// UpdateSvc 更新service
	UpdateSvc(*model.Svc) error

	// FindSvcByID 根据ID查找service
	FindSvcByID(int64) (*model.Svc, error)

	// FindAllSvc 查找全部service
	FindAllSvc() ([]model.Svc, error)

	// CreateSvcToK8s 创建服务到k8s
	CreateSvcToK8s(*svc.SvcInfo) error

	// UpdateSvcToK8s 更新服务到k8s
	UpdateSvcToK8s(*svc.SvcInfo) error

	// DeleteFromK8s 从k8s删除服务
	DeleteFromK8s(*model.Svc) error
}

// NewService 初始化Service
func NewService(serviceRepository repository.SvcRepository, clientSet *kubernetes.Clientset) SvcService {
	return &SvcDataService{
		ServiceRepository: serviceRepository,
		K8sClientSet:      clientSet,
	}
}

// SvcDataService service数据服务
type SvcDataService struct {
	// ServiceRepository 操作数据库接口
	ServiceRepository repository.SvcRepository

	// K8sClientSet k8s客户端集合
	K8sClientSet *kubernetes.Clientset
}

// AddSvc 添加service
func (s *SvcDataService) AddSvc(m *model.Svc) (int64, error) {
	return s.ServiceRepository.CreateSvc(m)
}

// DeleteSvc 删除service
func (s *SvcDataService) DeleteSvc(i int64) error {
	return s.ServiceRepository.DeleteSvcByID(i)
}

// UpdateSvc 更新service
func (s *SvcDataService) UpdateSvc(m *model.Svc) error {
	return s.ServiceRepository.UpdateSvc(m)
}

// FindSvcByID 根据ID查找service
func (s *SvcDataService) FindSvcByID(i int64) (*model.Svc, error) {
	return s.ServiceRepository.FindSvcByID(i)
}

// FindAllSvc 查找全部service
func (s *SvcDataService) FindAllSvc() ([]model.Svc, error) {
	return s.ServiceRepository.FindAll()
}

// CreateSvcToK8s 创建服务到k8s
func (s *SvcDataService) CreateSvcToK8s(info *svc.SvcInfo) error {
	svc := s.setService(info)

	// 查找是否存在指定服务
	_, err := s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Get(context.TODO(), info.SvcName, v12.GetOptions{})
	if err != nil {
		// 如果不存在 -> 创建service
		_, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Create(context.TODO(), svc, v12.CreateOptions{})
		if err != nil {
			// 创建失败
			common.Error(err)
			return err
		}

		// 创建成功
		return nil
	}

	// 说明已经存在
	common.Error("SvcService: " + info.SvcName + "已经存在")
	return errors.New("SvcService: " + info.SvcName + "已经存在")
}

// UpdateSvcToK8s 更新服务到k8s
func (s *SvcDataService) UpdateSvcToK8s(info *svc.SvcInfo) error {
	svc := s.setService(info)

	// 查找是否存在指定服务
	_, err := s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Get(context.TODO(), info.SvcName, v12.GetOptions{})
	if err != nil {
		// 不存在服务
		common.Error(err)
		return errors.New("SvcService：" + info.SvcName + "不存在服务请先创建")
	}

	// 之前存在 -> 可以更新
	_, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Update(context.TODO(), svc, v12.UpdateOptions{})
	if err != nil {
		// 更新失败
		common.Error(err)
		return err
	}
	// 更新成功
	common.Info("SvcService: ")
	return nil
}

// DeleteFromK8s 从k8s删除服务
func (s *SvcDataService) DeleteFromK8s(m *model.Svc) error {
	err := s.K8sClientSet.CoreV1().Services(m.SvcNamespace).Delete(context.TODO(), m.SvcName, v12.DeleteOptions{})
	if err != nil {
		// 删除失败
		common.Error(err)
		return err
	}

	// k8s成功删除,删除数据库中数据
	err = s.ServiceRepository.DeleteSvcByID(m.ID)
	if err != nil {
		common.Error(err)
		return err
	}
	common.Info("SvcService: " + strconv.FormatInt(m.ID, 10) + "删除成功")

	return nil
}

// setService 组装service信息
func (s *SvcDataService) setService(info *svc.SvcInfo) *v1.Service {
	svc := &v1.Service{}

	// 设置服务类型
	svc.TypeMeta = v12.TypeMeta{
		Kind:       "v1",
		APIVersion: "SvcService",
	}
	// 设置基础信息
	svc.ObjectMeta = v12.ObjectMeta{
		Name:         info.SvcName,
		GenerateName: info.SvcNamespace,
		Labels: map[string]string{
			"app-name": info.SvcPodName,
		},
		Annotations: map[string]string{
			"k8s/generated-by-zhao": "备注声明",
		},
	}
	// 设置服务的spec信息，采用ClusterIP模式
	svc.Spec = v1.ServiceSpec{
		Ports: s.getServicePort(info),
		Selector: map[string]string{
			"app-name": info.SvcPodName,
		},
		Type: "ClusterIP",
	}

	return svc
}

// getServicePort 获取服务端口
func (s *SvcDataService) getServicePort(info *svc.SvcInfo) []v1.ServicePort {
	var servicePort []v1.ServicePort

	for _, port := range info.SvcPort {
		// 将servicePort信息添加
		servicePort = append(servicePort, v1.ServicePort{
			Name:       "port-" + strconv.FormatInt(int64(port.SvcPort), 10),
			Protocol:   v1.Protocol(port.SvcPortProtocol),
			Port:       port.SvcPort,
			TargetPort: intstr.FromInt(int(port.SvcTargetPort)),
		})
	}

	return servicePort
}

package service

import (
	"context"
	"errors"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/networking/v1"
	v14 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"tini-paas/internal/route/model"
	"tini-paas/internal/route/proto/route"
	"tini-paas/internal/route/repository"
	"tini-paas/pkg/common"
)

// RouteService Route服务
type RouteService interface {
	// AddRoute 添加Route
	AddRoute(*model.Route) (int64, error)

	// DeleteRouteByID 删除Route
	DeleteRouteByID(int64) error

	// UpdateRoute 更新Route
	UpdateRoute(*model.Route) error

	// FindRouteByID 根据ID查找Route
	FindRouteByID(int64) (*model.Route, error)

	// FindAllRoute 查找全部Route
	FindAllRoute() ([]model.Route, error)

	// CreateRouteToK8s 创建Route到K8s
	CreateRouteToK8s(*route.RouteInfo) error

	// UpdateRouteToK8s 更新Route到k8s
	UpdateRouteToK8s(*route.RouteInfo) error

	// DeleteRouteFromK8s 从k8s删除Route
	DeleteRouteFromK8s(*model.Route) error
}

// NewRouteService 初始化route接口服务
func NewRouteService(routerRepository repository.RouteRepository, clientSet *kubernetes.Clientset) RouteService {
	return &RouteDataService{
		RouteRepository: routerRepository,
		K8sClientSet:    clientSet,
		deployment:      &v1.Deployment{},
	}
}

// RouteDataService route数据服务
type RouteDataService struct {
	// RouteRepository  操作数据库接口
	RouteRepository repository.RouteRepository

	// K8sClientSet k8s客户端集合
	K8sClientSet *kubernetes.Clientset

	// deployment 发布控制器
	deployment *v1.Deployment
}

// AddRoute 添加Route
func (r *RouteDataService) AddRoute(m *model.Route) (int64, error) {
	return r.RouteRepository.CreateRoute(m)
}

// DeleteRouteByID 删除Route
func (r *RouteDataService) DeleteRouteByID(i int64) error {
	return r.RouteRepository.DeleteRouteByID(i)
}

// UpdateRoute 更新Route
func (r *RouteDataService) UpdateRoute(m *model.Route) error {
	return r.RouteRepository.UpdateRoute(m)
}

// FindRouteByID 根据ID查找Route
func (r *RouteDataService) FindRouteByID(i int64) (*model.Route, error) {
	return r.RouteRepository.FindRouteByID(i)
}

// FindAllRoute 查找全部Route
func (r *RouteDataService) FindAllRoute() ([]model.Route, error) {
	return r.RouteRepository.FindAll()
}

// CreateRouteToK8s 创建Route到K8s
func (r *RouteDataService) CreateRouteToK8s(info *route.RouteInfo) error {
	// 组装信息
	ingress := r.setIngress(info)

	// 先查询之前是否存在
	_, err := r.K8sClientSet.NetworkingV1().Ingresses(info.RouteNamespace).Get(context.TODO(), info.RouteName, v14.GetOptions{})
	if err != nil {
		// 之前不存在 -> 可以创建
		_, err = r.K8sClientSet.NetworkingV1().Ingresses(info.RouteNamespace).Create(context.TODO(), ingress, v14.CreateOptions{})
		if err != nil {
			// 创建失败
			common.Error(err)
			return err
		}
		return nil
	}

	// 之前存在ingress
	// TODO
	common.Error("路由：" + info.RouteName + "已经存在")
	return errors.New("路由：" + info.RouteName + "已经存在")
}

// UpdateRouteToK8s 更新Route到k8s
func (r *RouteDataService) UpdateRouteToK8s(info *route.RouteInfo) error {
	ingress := r.setIngress(info)

	_, err := r.K8sClientSet.NetworkingV1().Ingresses(info.RouteNamespace).Update(context.TODO(), ingress, v14.UpdateOptions{})
	if err != nil {
		// 更新失败
		common.Error(err)
		return err
	}
	return nil
}

// DeleteRouteFromK8s 从k8s删除Route
func (r *RouteDataService) DeleteRouteFromK8s(m *model.Route) error {
	// 先删除k8s
	err := r.K8sClientSet.NetworkingV1().Ingresses(m.RouteNamespace).Delete(context.TODO(), m.RouteName, v14.DeleteOptions{})
	if err != nil {
		// 删除失败
		common.Error(err)
		return err
	}

	// 删除k8s成功后，删除数据库
	err = r.RouteRepository.DeleteRouteByID(m.ID)
	if err != nil {
		common.Error(err)
		return err
	}

	common.Info("删除 ingress ID: " + strconv.FormatInt(m.ID, 10) + "成功")
	return nil
}

// setIngress 封装ingress
func (r *RouteDataService) setIngress(info *route.RouteInfo) *v12.Ingress {
	router := &v12.Ingress{}

	// 设置路由
	router.TypeMeta = v14.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "v1",
	}

	// 设置路由基础信息
	router.ObjectMeta = v14.ObjectMeta{
		Name:      info.RouteName,
		Namespace: info.RouteNamespace,
		Labels: map[string]string{
			"app-name": info.RouteName,
			"author":   "router",
		},
		Annotations: map[string]string{
			"k8s/generated": "paasmicro",
		},
	}

	// 使用 ingress-nginx
	className := "nginx"
	// 设置路由 spec 信息
	router.Spec = v12.IngressSpec{
		IngressClassName: &className,
		DefaultBackend:   nil, // 默认访问服务
		TLS:              nil, // https
		Rules:            r.getIngressPath(info),
	}

	return router
}

// getIngressPath 封装Ingress路径
func (r *RouteDataService) getIngressPath(info *route.RouteInfo) []v12.IngressRule {
	var path []v12.IngressRule

	// 1、设置host
	pathRule := v12.IngressRule{
		Host: info.RouteHost,
	}

	// 2、设置path
	var ingressPath []v12.HTTPIngressPath
	for _, routePath := range info.RoutePath {
		pathType := v12.PathTypePrefix

		// 将信息写入
		ingressPath = append(ingressPath, v12.HTTPIngressPath{
			Path:     routePath.RoutePathName,
			PathType: &pathType,
			Backend: v12.IngressBackend{
				Service: &v12.IngressServiceBackend{
					Name: routePath.RouteBackendService,
					Port: v12.ServiceBackendPort{
						Number: routePath.RouteBackendServicePort,
					},
				},
			},
		})
	}

	// 3、赋值path
	pathRule.IngressRuleValue = v12.IngressRuleValue{
		HTTP: &v12.HTTPIngressRuleValue{
			Paths: ingressPath,
		},
	}
	path = append(path, pathRule)
	return path
}

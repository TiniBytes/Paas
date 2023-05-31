package handler

import (
	"context"
	"tini-paas/internal/route/model"
	"tini-paas/internal/route/proto/route"
	"tini-paas/internal/route/service"
	"tini-paas/pkg/common"
)

// RouteHandler 调用底层, 实现服务层接口
type RouteHandler struct {
	RouteService service.RouteService
}

func (r *RouteHandler) AddRoute(ctx context.Context, info *route.RouteInfo, response *route.Response) error {
	routeModel := &model.Route{}

	// 通过json tag 将info映射到routeModel上
	err := common.SwapTo(info, routeModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 创建route到k8s
	err = r.RouteService.CreateRouteToK8s(info)
	if err != nil {
		// 创建失败
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 创建成功 -> 写入数据库
	_, err = r.RouteService.AddRoute(routeModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	return nil
}

func (r *RouteHandler) DeleteRoute(ctx context.Context, id *route.RouteID, response *route.Response) error {
	// 先查找是否存在
	routeModel, err := r.RouteService.FindRouteByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 从k8s中删除 -> 在数据库中删除
	err = r.RouteService.DeleteRouteFromK8s(routeModel)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

func (r *RouteHandler) UpdateRoute(ctx context.Context, info *route.RouteInfo, response *route.Response) error {
	// 先更新k8s
	err := r.RouteService.UpdateRouteToK8s(info)
	if err != nil {
		common.Error(err)
		return err
	}

	// 查询数据库信息
	routeModel, err := r.RouteService.FindRouteByID(info.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 更新数据库
	err = common.SwapTo(info, routeModel)
	if err != nil {
		common.Error(err)
		return err
	}
	return r.RouteService.UpdateRoute(routeModel)
}

func (r *RouteHandler) FindRouteByID(ctx context.Context, id *route.RouteID, info *route.RouteInfo) error {
	routeModel, err := r.RouteService.FindRouteByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据转换
	err = common.SwapTo(routeModel, info)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

func (r *RouteHandler) FindAllRoute(ctx context.Context, all *route.FindAll, rsp *route.AllRoute) error {
	allRoute, err := r.RouteService.FindAllRoute()
	if err != nil {
		common.Error(err)
		return err
	}

	// 整理格式
	for _, v := range allRoute {
		// 创建实例
		routeInfo := &route.RouteInfo{}

		// 将查询结果转换
		err = common.SwapTo(v, routeInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		rsp.RouteInfo = append(rsp.RouteInfo, routeInfo)
	}
	return nil
}

package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"tini-paas/api/routeapi/proto/routeApi"
	"tini-paas/internal/route/proto/route"
	"tini-paas/pkg/common"
	"tini-paas/plugin/form"
)

// RouteApi Api接口
type RouteApi struct {
	RouteService route.RouteService
}

// AddRoute routeApi.AddRoute 通过API向外暴露为/routeApi/AddRoute，接收http请求
// 即：/routeApi/AddRoute 请求会调用go.micro.api.AddRoute 服务的routeApi.AddRoute 方法
func (r *RouteApi) AddRoute(ctx context.Context, req *routeApi.Request, rsp *routeApi.Response) error {
	addRouteInfo := &route.RouteInfo{}
	routePathName, ok := req.Post["route_path_name"]

	if ok && len(routePathName.Values) > 0 {
		// port
		port, err := strconv.ParseInt(req.Post["route_backend_service_port"].Values[0], 10, 64)
		if err != nil {
			common.Error(err)
			return err
		}

		// 处理Path路径
		routePath := &route.RoutePath{
			RoutePathName:           req.Post["route_path_name"].Values[0],
			RouteBackendService:     req.Post["route_backend_service"].Values[0],
			RouteBackendServicePort: int32(port),
		}

		// 组装信息
		addRouteInfo.RoutePath = append(addRouteInfo.RoutePath, routePath)
	}

	// 执行添加
	form.FormToRouteStruct(req.Post, addRouteInfo)
	response, err := r.RouteService.AddRoute(ctx, addRouteInfo)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(response)
	rsp.Body = string(bytes)
	return nil
}

// DeleteRoute routeApi.DeleteRoute 通过API向外暴露为/routeApi/DeleteRoute，接收http请求
// 即：/routeApi/DeleteRoute 请求会调用go.micro.api.DeleteRoute 服务的routeApi.DeleteRoute 方法
func (r *RouteApi) DeleteRoute(ctx context.Context, req *routeApi.Request, rsp *routeApi.Response) error {
	// 先查询id是否存在
	if _, ok := req.Get["route_id"]; !ok {
		rsp.StatusCode = 200
		return errors.New("参数异常")
	}

	// 获取RouteID
	routeIDString := req.Get["route_id"].Values[0]
	routeID, err := strconv.ParseInt(routeIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 执行删除
	response, err := r.RouteService.DeleteRoute(ctx, &route.RouteID{
		Id: routeID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(response)
	rsp.Body = string(bytes)
	return nil
}

// UpdateRoute routeApi.UpdateRoute 通过API向外暴露为/routeApi/UpdateRoute，接收http请求
// 即：/routeApi/UpdateRoute 请求会调用go.micro.api.UpdateRoute 服务的routeApi.UpdateRoute 方法
func (r *RouteApi) UpdateRoute(ctx context.Context, req *routeApi.Request, rsp *routeApi.Response) error {
	//TODO implement me
	panic("implement me")
}

// FindRouteByID routeApi.FindRouteByID 通过API向外暴露为/routeApi/FindRouteByID，接收http请求
// 即：/routeApi/FindRouteByID 请求会调用go.micro.api.FindRouteByID 服务的routeApi.FindRouteByID 方法
func (r *RouteApi) FindRouteByID(ctx context.Context, req *routeApi.Request, rsp *routeApi.Response) error {
	// 查询id是否存在
	if _, ok := req.Get["route_id"]; !ok {
		rsp.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取RouteID
	routeIDString := req.Get["route_id"].Values[0]
	routeID, err := strconv.ParseInt(routeIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 执行查询
	routeInfo, err := r.RouteService.FindRouteByID(ctx, &route.RouteID{
		Id: routeID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 回写消息
	bytes, _ := json.Marshal(routeInfo)
	rsp.StatusCode = 200
	rsp.Body = string(bytes)
	return nil
}

// Call routeApi.Call 通过API向外暴露为/routeApi/Call，接收http请求
// 即：/routeApi/Call 请求会调用go.micro.api.Call 服务的routeApi.Call 方法
func (r *RouteApi) Call(ctx context.Context, req *routeApi.Request, rsp *routeApi.Response) error {
	allRoute, err := r.RouteService.FindAllRoute(ctx, &route.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(allRoute)
	rsp.Body = string(bytes)
	return nil
}

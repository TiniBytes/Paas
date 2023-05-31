package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"tini-paas/api/svcapi/proto/svcApi"
	"tini-paas/internal/svc/proto/svc"
	"tini-paas/pkg/common"
	"tini-paas/plugin/form"
)

// SvcApi Svc
type SvcApi struct {
	SvcService svc.SvcService
}

// AddSvc svcApi.AddSvc 通过API向外暴露为/svcApi/AddSvc，接收http请求
// 即：/svcApi/AddSvc 请求会调用go.micro.api.svcApi 服务的svcApi.AddSvc 方法
func (s *SvcApi) AddSvc(ctx context.Context, req *svcApi.Request, rsp *svcApi.Response) error {
	fmt.Println("接收到 svcApi.AddSvc 的请求")

	// 处理port
	addSvcInfo := &svc.SvcInfo{}
	svcType, ok := req.Post["svc_type"]
	if ok && len(svcType.Values) > 0 {
		// 组装svcPort信息
		svcPort := &svc.SvcPort{}

		// 选择四种Service类型
		switch svcType.Values[0] {
		case "ClusterIP":
			// 解析服务端口
			port, err := strconv.ParseInt(req.Post["svc_port"].Values[0], 10, 32)
			if err != nil {
				common.Error(err)
				return err
			}
			svcPort.SvcPort = int32(port)

			// 解析目标端口
			targetPort, err := strconv.ParseInt(req.Post["svc_target_port"].Values[0], 10, 64)
			if err != nil {
				common.Error(err)
				return err
			}
			svcPort.SvcTargetPort = int32(targetPort)

			// 解析端口协议
			svcPort.SvcPortProtocol = req.Post["svc_port_protocol"].Values[0]

			// 将本组消息写入serviceInfo
			addSvcInfo.SvcPort = append(addSvcInfo.SvcPort, svcPort)

			// TODO 其它service类型
		default:
			return errors.New("类型不支持")
		}
	}

	// 将form表单映射到结构体中
	form.FormToSvcStruct(req.Post, addSvcInfo)

	// 添加svc
	response, err := s.SvcService.AddSvc(ctx, addSvcInfo)
	if err != nil {
		common.Error(err)
		return err
	}

	// 状态回写
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(response)
	rsp.Body = string(bytes)
	return nil
}

// DeleteSvcByID svcApi.DeleteSvcByID 通过API向外暴露为/svcApi/DeleteSvcByID，接收http请求
// 即：/svcApi/DeleteSvcByID 请求会调用go.micro.api.DeleteSvcByID 服务的svcApi.DeleteSvcByID 方法
func (s *SvcApi) DeleteSvcByID(ctx context.Context, req *svcApi.Request, rsp *svcApi.Response) error {
	fmt.Println("接收到 svcApi.DeleteSvcByID 的请求")

	// 较验参数
	if _, ok := req.Get["svc_id"]; !ok {
		return errors.New("参数异常")
	}

	// 获取要删除的ID
	svcIDString := req.Get["svc_id"].Values[0]
	ID, err := strconv.ParseInt(svcIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务删除
	response, err := s.SvcService.DeleteSvc(ctx, &svc.SvcID{
		Id: ID,
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

// UpdateSvc svcApi.UpdateSvc 通过API向外暴露为/svcApi/UpdateSvc，接收http请求
// 即：/svcApi/UpdateSvc 请求会调用go.micro.api.UpdateSvc 服务的svcApi.UpdateSvc 方法
func (s *SvcApi) UpdateSvc(ctx context.Context, request *svcApi.Request, response *svcApi.Response) error {
	//TODO implement me
	panic("implement me")
}

// FindSvcByID svcApi.FindSvcByID 通过API向外暴露为/svcApi/FindSvcByID，接收http请求
// 即：/svcApi/FindSvcByID 请求会调用go.micro.api.FindSvcByID 服务的svcApi.FindSvcByID 方法
func (s *SvcApi) FindSvcByID(ctx context.Context, request *svcApi.Request, response *svcApi.Response) error {
	fmt.Println("接收到 svcApi.FindSvcByID 的请求")

	// 较验请求中的pod_id
	if _, ok := request.Get["svc_id"]; !ok {
		response.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取svc相关信息
	svcIDString := request.Get["svc_id"].Values[0]
	ID, err := strconv.ParseInt(svcIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 访问后端的svcService
	svcInfo, err := s.SvcService.FindSvcByID(ctx, &svc.SvcID{
		Id: ID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 将info信息写入到response
	response.StatusCode = 200
	bytes, _ := json.Marshal(svcInfo)
	response.Body = string(bytes)
	return nil
}

// Call svcApi.Call 通过API向外暴露为/svcApi/Call，接收http请求
// 即：/svcApi/Call 请求会调用go.micro.api.Call 服务的svcApi.Call 方法
func (s *SvcApi) Call(ctx context.Context, req *svcApi.Request, rsp *svcApi.Response) error {
	fmt.Println("接收到 svcApi.Call 的请求")

	allSvc, err := s.SvcService.FindAllSvc(ctx, &svc.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(allSvc)
	rsp.Body = string(bytes)
	return nil
}

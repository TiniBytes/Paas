package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"tini-paas/api/podapi/proto/podApi"
	"tini-paas/internal/pod/proto/pod"
	"tini-paas/pkg/common"
	"tini-paas/plugin/form"
)

// PodApi pod
type PodApi struct {
	PodService pod.PodService
}

// FindPodByID 查找pod
// podApi.FindPodByID 通过API向外暴露为/podApi/FindPodByID, 接收http请求
// 即：/podApi/FindPodByID 请求会调用go.micro.api.podApi 服务的podApi.FindPodByID方法
func (p *PodApi) FindPodByID(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	fmt.Println("接收到 podApi.FindPodByID 的请求")

	// 拿到请求中的pod_id
	if _, ok := req.Get["pod_id"]; !ok {
		rsp.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取pod相关信息
	podIDString := req.Get["pod_id"].Values[0]
	ID, err := strconv.ParseInt(podIDString, 10, 64)
	if err != nil {
		return err
	}

	// 访问后端的podService
	podInfo, err := p.PodService.FindPodByID(ctx, &pod.PodID{
		Id: ID,
	})
	if err != nil {
		return err
	}

	// 将info信息写入到response
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(podInfo)
	rsp.Body = string(bytes)
	return nil
}

// AddPod 添加pod
// PodApi.AddPod 通过API向外暴露为/podApi/AddPod, 接收http请求
// 即：/podApi/AddPod 请求会调用go.micro.api.PodApi 服务的PodApi.AddPod方法
func (p *PodApi) AddPod(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	fmt.Println("接收到 podApi.AddPod 的请求")
	addPodInfo := &pod.PodInfo{}

	// 处理port
	dataSlice, ok := req.Post["pod_port"]
	if ok {
		// 特殊处理
		var podSlice []*pod.PodPort
		for _, v := range dataSlice.Values {
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				common.Error(err)
			}

			// 封装port
			port := &pod.PodPort{
				ContainerPort: int32(i),
				Protocol:      "TCP",
			}
			podSlice = append(podSlice, port)
		}
		// 信息写入
		addPodInfo.PodPort = podSlice
	}

	// 将form表单映射到结构体中
	form.FromToPodStruct(req.Post, addPodInfo)

	// 添加pod
	response, err := p.PodService.AddPod(ctx, addPodInfo)
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

// DeletePodByID 删除pod
// PodApi.DeletedPodByID 通过API向外暴露为/podApi/DeletePodByID, 接收http请求
// 即：/podApi/DeletePodByID 请求会调用go.micro.api.PodApi 服务的PodApi.DeletePodByID方法
func (p *PodApi) DeletePodByID(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	fmt.Println("接收到 podApi.DeletePodByID 的请求")
	if _, ok := req.Get["pod_id"]; !ok {
		return errors.New("参数异常")
	}

	// 获取要删除的podID
	podIDString := req.Get["pod_id"].Values[0]
	podID, err := strconv.ParseInt(podIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用服务，删除pod
	response, err := p.PodService.DeletePod(ctx, &pod.PodID{
		Id: podID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 回写数据
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(response)
	rsp.Body = string(bytes)
	return nil
}

// UpdatePod 更新pod
// PodApi.UpdatePod 通过API向外暴露为/podApi/UpdatePod, 接收http请求
// 即：/podApi/UpdatePod 请求会调用go.micro.api.PodApi 服务的PodApi.UpdatePod方法
func (p *PodApi) UpdatePod(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {

	return nil
}

// Call 默认方法
// PodApi.Call 通过API向外暴露为/podApi/Call, 接收http请求
// 即：/podApi/Call 请求会调用go.micro.api.Call 服务的PodApi.Call方法
func (p *PodApi) Call(ctx context.Context, req *podApi.Request, rsp *podApi.Response) error {
	fmt.Println("接收到 podApi.Call 的请求")
	allPod, err := p.PodService.FindAllPod(ctx, &pod.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}

	// 回写数据
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(allPod)
	rsp.Body = string(bytes)
	return nil
}

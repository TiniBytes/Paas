package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"tini-paas/api/volumeapi/proto/volumeApi"
	"tini-paas/internal/volume/proto/volume"
	"tini-paas/pkg/common"
	"tini-paas/plugin/form"
)

// VolumeApi handler 调用volume的客户端API接口
type VolumeApi struct {
	VolumeServer volume.VolumeService
}

func (v *VolumeApi) AddVolume(ctx context.Context, req *volumeApi.Request, rsp *volumeApi.Response) error {
	addVolumeInfo := &volume.VolumeInfo{}

	// 将req.Post信息转换为VolumeInfo
	form.FormToVolumeStruct(req.Post, addVolumeInfo)

	// 添加volume
	response, err := v.VolumeServer.AddVolume(ctx, addVolumeInfo)
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

func (v *VolumeApi) DeleteVolume(ctx context.Context, req *volumeApi.Request, rsp *volumeApi.Response) error {
	// 先查询id是否可以解析
	if _, ok := req.Get["volume_id"]; !ok {
		rsp.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取volume_id
	volumeIDString := req.Get["volume_id"].Values[0]
	volumeID, err := strconv.ParseInt(volumeIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 执行删除服务
	response, err := v.VolumeServer.DeleteVolume(ctx, &volume.VolumeID{
		Id: volumeID,
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

func (v *VolumeApi) UpdateVolume(ctx context.Context, req *volumeApi.Request, rsp *volumeApi.Response) error {
	//TODO implement me
	panic("implement me")
}

func (v *VolumeApi) FindVolumeByID(ctx context.Context, req *volumeApi.Request, rsp *volumeApi.Response) error {
	if _, ok := req.Get["volume_id"]; !ok {
		rsp.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取volume_id
	volumeIDString := req.Get["volume_id"].Values[0]
	volumeID, err := strconv.ParseInt(volumeIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 执行查询
	volumeInfo, err := v.VolumeServer.FindVolumeByID(ctx, &volume.VolumeID{
		Id: volumeID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(volumeInfo)
	rsp.Body = string(bytes)
	return nil
}

func (v *VolumeApi) Call(ctx context.Context, req *volumeApi.Request, rsp *volumeApi.Response) error {
	allVolume, err := v.VolumeServer.FindAllVolume(ctx, &volume.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	rsp.StatusCode = 200
	bytes, _ := json.Marshal(allVolume)
	rsp.Body = string(bytes)
	return nil
}

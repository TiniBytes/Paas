package handler

import (
	"context"
	"strconv"
	"tini-paas/internal/volume/model"
	"tini-paas/internal/volume/proto/volume"
	"tini-paas/internal/volume/service"
	"tini-paas/pkg/common"
)

// VolumeHandler 操作接口
type VolumeHandler struct {
	VolumeService service.VolumeService
}

func (v *VolumeHandler) AddVolume(ctx context.Context, info *volume.VolumeInfo, response *volume.Response) error {
	volume := &model.Volume{}

	// 将info信息映射到volume
	err := common.SwapTo(info, volume)
	if err != nil {
		common.Error(err)
		return err
	}

	// 创建volume, 先在k8s中创建
	err = v.VolumeService.CreateVolumeToK8s(info)
	if err != nil {
		common.Error(err)
		return err
	}

	// 创建成功后写入数据库
	volumeID, err := v.VolumeService.AddVolume(volume)
	if err != nil {
		common.Error(err)
		return err
	}

	// 回写数据
	common.Info("Svc 添加成功，ID为：" + strconv.FormatInt(volumeID, 10))
	response.Msg = "Svc 添加成功，ID为：" + strconv.FormatInt(volumeID, 10)
	return nil
}

func (v *VolumeHandler) DeleteVolume(ctx context.Context, id *volume.VolumeID, response *volume.Response) error {
	// 先查询是否存在
	volumeModel, err := v.VolumeService.FindVolume(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 在k8s中删除
	err = v.VolumeService.DeleteVolumeFromK8s(volumeModel)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

func (v *VolumeHandler) UpdateVolume(ctx context.Context, info *volume.VolumeInfo, response *volume.Response) error {
	//TODO implement me
	panic("implement me")
}

func (v *VolumeHandler) FindVolumeByID(ctx context.Context, id *volume.VolumeID, info *volume.VolumeInfo) error {
	volumeModel, err := v.VolumeService.FindVolume(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	err = common.SwapTo(volumeModel, info)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

func (v *VolumeHandler) FindAllVolume(ctx context.Context, all *volume.FindAll, allVolume *volume.AllVolume) error {
	volumes, err := v.VolumeService.FindAllVolume()
	if err != nil {
		common.Error(err)
		return err
	}

	// 整理格式
	for _, vol := range volumes {
		// 创建实例
		volumeInfo := &volume.VolumeInfo{}
		// 数据转换
		err = common.SwapTo(vol, volumeInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		// 数据合并
		allVolume.VolumeInfo = append(allVolume.VolumeInfo, volumeInfo)
	}
	return nil
}

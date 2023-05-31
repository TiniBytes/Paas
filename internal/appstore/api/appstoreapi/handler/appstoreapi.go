package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"tini-paas/api/appstoreapi/proto/appstoreApi"
	"tini-paas/internal/appstore/proto/appstore"
	"tini-paas/pkg/common"
	"tini-paas/plugin/form"
)

// AppStoreApi 中间件API处理（对API后端接口方法的实现）
type AppStoreApi struct {
	AppStoreService appstore.AppStoreService
}

func (a *AppStoreApi) AddAppStore(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	addAppStore := &appstore.AppStoreInfo{}

	// form表单数据映射
	form.FormToAppStoreStruct(request.Post, addAppStore)

	// 设置图片
	a.setImage(request, addAppStore)
	// 设置pod
	a.setPod(request, addAppStore)
	// 设置中间件模板
	a.setMiddle(request, addAppStore)
	// 设置存储
	a.setVolume(request, addAppStore)

	// 调用后端服务添加
	rsp, err := a.AppStoreService.AddAppStore(ctx, addAppStore)
	if err != nil {
		common.Error(err)
		return err
	}

	// 状态回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

func (a *AppStoreApi) DeleteAppStore(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	id, err := a.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行删除
	rsp, err := a.AppStoreService.DeleteAppStore(ctx, &appstore.AppStoreID{
		Id: id,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 状态回写
	response.StatusCode = 200
	bytes, err := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

func (a *AppStoreApi) UpdateAppStore(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	id, err := a.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 查询info信息
	info, err := a.AppStoreService.FindAppStoreByID(ctx, &appstore.AppStoreID{
		Id: id,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// form表单数据映射
	form.FormToAppStoreStruct(request.Post, info)

	// 调用后端服务执行更新
	rsp, err := a.AppStoreService.UpdateAppStore(ctx, info)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	rsp.Msg = string(bytes)
	return nil
}

func (a *AppStoreApi) FindAppStoreByID(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	id, err := a.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := a.AppStoreService.FindAppStoreByID(ctx, &appstore.AppStoreID{
		Id: id,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

func (a *AppStoreApi) Call(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	allAppStore, err := a.AppStoreService.FindAllAppStore(ctx, &appstore.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(allAppStore)
	response.Body = string(bytes)
	return nil
}

func (a *AppStoreApi) AddInstallNum(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	id, err := a.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := a.AppStoreService.AddInstallNum(ctx, &appstore.AppStoreID{
		Id: id,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

func (a *AppStoreApi) GetInstallNum(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	id, err := a.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := a.AppStoreService.GetInstallNum(ctx, &appstore.AppStoreID{
		Id: id,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

func (a *AppStoreApi) AddViewNum(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	id, err := a.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := a.AppStoreService.AddViewNum(ctx, &appstore.AppStoreID{
		Id: id,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

func (a *AppStoreApi) GetViewNum(ctx context.Context, request *appstoreApi.Request, response *appstoreApi.Response) error {
	id, err := a.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := a.AppStoreService.GetViewNum(ctx, &appstore.AppStoreID{
		Id: id,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// setImage 设置图片
func (a *AppStoreApi) setImage(request *appstoreApi.Request, info *appstore.AppStoreInfo) {
	data, ok := request.Post["app_image"]
	if !ok {
		return
	}

	var imageSlice []*appstore.AppImage
	for _, value := range data.Values {
		img := &appstore.AppImage{
			AppImageSrc: value,
		}

		imageSlice = append(imageSlice, img)
	}

	// 写入info
	info.AppImage = imageSlice
}

// setPod 设置pod
func (a *AppStoreApi) setPod(request *appstoreApi.Request, info *appstore.AppStoreInfo) {
	data, ok := request.Post["app_pod"]
	if !ok {
		return
	}

	var podSlice []*appstore.AppPod
	for _, value := range data.Values {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			common.Error(err)
			continue
		}

		pod := &appstore.AppPod{
			AppPodId: id,
		}

		podSlice = append(podSlice, pod)
	}

	// 写入info
	info.AppPod = podSlice
}

// setMiddle 设置中间件模板
func (a *AppStoreApi) setMiddle(request *appstoreApi.Request, info *appstore.AppStoreInfo) {
	data, ok := request.Post["app_middle"]
	if !ok {
		return
	}

	var middleSlice []*appstore.AppMiddle
	for _, value := range data.Values {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			common.Error(err)
			continue
		}

		middle := &appstore.AppMiddle{
			AppMiddleId: id,
		}

		middleSlice = append(middleSlice, middle)
	}

	// 写入info
	info.AppMiddle = middleSlice
}

// setVolume 设置存储
func (a *AppStoreApi) setVolume(request *appstoreApi.Request, info *appstore.AppStoreInfo) {
	data, ok := request.Post["app_volume"]
	if !ok {
		return
	}

	var volumeSlice []*appstore.AppVolume
	for _, value := range data.Values {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			common.Error(err)
			continue
		}

		volume := &appstore.AppVolume{
			AppVolumeId: id,
		}

		volumeSlice = append(volumeSlice, volume)
	}

	// 写入info
	info.AppVolume = volumeSlice
}

// getID 获取ID
func (a *AppStoreApi) getID(request *appstoreApi.Request) (int64, error) {
	// 检验
	if _, ok := request.Get["app_id"]; !ok {
		return 0, errors.New("参数异常")
	}

	// 获取ID进行转化
	idString := request.Get["app_id"].Values[0]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		common.Error(err)
		return 0, err
	}

	return id, nil
}

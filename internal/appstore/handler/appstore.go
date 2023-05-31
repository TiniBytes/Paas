package handler

import (
	"context"
	"strconv"
	"tini-paas/internal/appstore/model"
	"tini-paas/internal/appstore/proto/appstore"
	"tini-paas/internal/appstore/service"
	"tini-paas/pkg/common"
)

// AppStoreHandler 应用市场（API接口实现）
type AppStoreHandler struct {
	AppStoreDataService service.AppStoreService
}

func (a *AppStoreHandler) AddAppStore(ctx context.Context, info *appstore.AppStoreInfo, response *appstore.Response) error {
	appStoreModel := &model.AppStore{}

	// 将info数据映射到appStoreModel
	err := common.SwapTo(info, appStoreModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 执行添加
	appStoreID, err := a.AppStoreDataService.AddAppStore(appStoreModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.Msg = "应用市场应用添加成功" + strconv.FormatInt(appStoreID, 10)
	common.Info(response.Msg)
	return nil
}

func (a *AppStoreHandler) DeleteAppStore(ctx context.Context, id *appstore.AppStoreID, response *appstore.Response) error {
	return a.AppStoreDataService.DeleteAppStore(id.Id)
}

func (a *AppStoreHandler) UpdateAppStore(ctx context.Context, info *appstore.AppStoreInfo, response *appstore.Response) error {
	// 先查询之前是否存在
	appStoreModel, err := a.AppStoreDataService.FindAppStoreByID(info.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将新数据写入appStoreModel
	err = common.SwapTo(info, appStoreModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 执行更新
	return a.AppStoreDataService.UpdateAppStore(appStoreModel)
}

func (a *AppStoreHandler) FindAppStoreByID(ctx context.Context, id *appstore.AppStoreID, info *appstore.AppStoreInfo) error {
	appStoreModel, err := a.AppStoreDataService.FindAppStoreByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据转换
	return common.SwapTo(appStoreModel, info)
}

func (a *AppStoreHandler) FindAllAppStore(ctx context.Context, all *appstore.FindAll, store *appstore.AllAppStore) error {
	allAppStore, err := a.AppStoreDataService.FindAllAppStore()
	if err != nil {
		common.Error(err)
		return err
	}

	for _, m := range allAppStore {
		appStoreInfo := &appstore.AppStoreInfo{}
		err = common.SwapTo(m, appStoreInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		store.AppStoreInfo = append(store.AppStoreInfo, appStoreInfo)
	}
	return nil
}

func (a *AppStoreHandler) AddInstallNum(ctx context.Context, id *appstore.AppStoreID, response *appstore.Response) error {
	err := a.AppStoreDataService.AddInstallNum(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	response.Msg = "ok"
	return nil
}

func (a *AppStoreHandler) GetInstallNum(ctx context.Context, id *appstore.AppStoreID, number *appstore.Number) error {
	number.Num = a.AppStoreDataService.GetInstallNum(id.Id)
	return nil
}

func (a *AppStoreHandler) AddViewNum(ctx context.Context, id *appstore.AppStoreID, response *appstore.Response) error {
	err := a.AppStoreDataService.AddViewNum(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	response.Msg = "ok"
	return nil
}

func (a *AppStoreHandler) GetViewNum(ctx context.Context, id *appstore.AppStoreID, number *appstore.Number) error {
	number.Num = a.AppStoreDataService.GetViewNum(id.Id)
	return nil
}

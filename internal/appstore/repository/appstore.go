package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/appstore/model"
	"tini-paas/pkg/common"
)

// AppStoreRepository 云应用商店接口
type AppStoreRepository interface {
	InitTable() error
	CreateAppStore(store *model.AppStore) (int64, error)
	DeleteAppStore(id int64) error
	UpdateAppStore(store *model.AppStore) error
	FindAppStoreByID(id int64) (*model.AppStore, error)
	FindAll() ([]model.AppStore, error)
	AddInstallNumber(id int64) error
	GetInstallNumber(id int64) int64
	AddViewNumber(id int64) error
	GetViewNumber(id int64) int64
}

// NewAppStoreRepository 初始化数据操作对象
func NewAppStoreRepository(db *gorm.DB) AppStoreRepository {
	return &AppStore{
		db: db,
	}
}

// AppStore 应用市场对象
type AppStore struct {
	db *gorm.DB
}

// InitTable 初始化表
func (a *AppStore) InitTable() error {
	return a.db.CreateTable(&model.AppStore{}, &model.AppCategory{}, &model.AppComment{}, &model.AppImage{}, &model.AppIsv{}, &model.AppMiddle{}, &model.AppPod{}, &model.AppVolume{}).Error
}

// CreateAppStore 创建应用市场
func (a *AppStore) CreateAppStore(store *model.AppStore) (int64, error) {
	return store.ID, a.db.Create(store).Error
}

// DeleteAppStore 删除应用市场
func (a *AppStore) DeleteAppStore(i int64) error {
	// 开启事务
	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			// 遇到问题回滚
			tx.Callback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除应用市场
	err := a.db.Delete(&model.AppStore{}).Where("id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除应用图片
	err = a.db.Delete(&model.AppImage{}).Where("app_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除Pod组合
	err = a.db.Delete(&model.AppPod{}).Where("app_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除中间件组合
	err = a.db.Delete(&model.AppMiddle{}).Where("app_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除存储组合
	err = a.db.Delete(&model.AppVolume{}).Where("app_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除评论
	err = a.db.Delete(&model.AppComment{}).Where("app_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}
	return tx.Commit().Error
}

// UpdateAppStore 更新应用市场
func (a *AppStore) UpdateAppStore(store *model.AppStore) error {
	return a.db.Model(store).Update(store).Error
}

// FindAppStoreByID 查询应用市场数据
func (a *AppStore) FindAppStoreByID(appStoreID int64) (*model.AppStore, error) {
	store := &model.AppStore{}
	return store, a.db.Preload("AppImage").Preload("AppPod").Preload("AppMiddle").Preload("AppVolume").Preload("AppComment").First(store, appStoreID).Error
}

// FindAll 查询全部应用市场数据
func (a *AppStore) FindAll() ([]model.AppStore, error) {
	var appStoreAll []model.AppStore
	return appStoreAll, a.db.Find(appStoreAll).Error
}

// AddInstallNumber 添加安装数量
func (a *AppStore) AddInstallNumber(i int64) error {
	return a.db.Model(&model.AppStore{}).Where("id = ?", i).UpdateColumn("app_install", gorm.Expr("app_install + ?", 1)).Error
}

// GetInstallNumber 获取安装数量
func (a *AppStore) GetInstallNumber(i int64) int64 {
	appStore, err := a.FindAppStoreByID(i)
	if err != nil {
		common.Error(err)
		return 0
	}
	return appStore.AppInstall
}

// AddViewNumber 添加浏览量
func (a *AppStore) AddViewNumber(i int64) error {
	return a.db.Model(&model.AppStore{}).Where("id = ?", i).UpdateColumn("app_views", gorm.Expr("app_views + ?", 1)).Error
}

// GetViewNumber 获取浏览量
func (a *AppStore) GetViewNumber(i int64) int64 {
	appStore, err := a.FindAppStoreByID(i)
	if err != nil {
		common.Error(err)
		return 0
	}
	return appStore.AppViews
}

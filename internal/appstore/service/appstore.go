package service

import (
	"k8s.io/client-go/kubernetes"
	"tini-paas/internal/appstore/model"
	"tini-paas/internal/appstore/repository"
)

// AppStoreService 应用市场服务接口
type AppStoreService interface {
	AddAppStore(store *model.AppStore) (int64, error)
	DeleteAppStore(id int64) error
	UpdateAppStore(store *model.AppStore) error
	FindAppStoreByID(id int64) (*model.AppStore, error)
	FindAllAppStore() ([]model.AppStore, error)

	AddInstallNum(id int64) error
	GetInstallNum(id int64) int64
	AddViewNum(id int64) error
	GetViewNum(id int64) int64
}

// NewAppStoreService 初始化应用市场服务
func NewAppStoreService(storeRepository repository.AppStoreRepository, client *kubernetes.Clientset) AppStoreService {
	return &AppStore{
		AppStoreRepository: storeRepository,
	}
}

// AppStore 应用市场服务对象
type AppStore struct {
	// AppStoreRepository 数据库对象
	AppStoreRepository repository.AppStoreRepository
}

// AddAppStore 添加应用市场
func (a *AppStore) AddAppStore(store *model.AppStore) (int64, error) {
	return a.AppStoreRepository.CreateAppStore(store)
}

// DeleteAppStore 删除应用市场
func (a *AppStore) DeleteAppStore(id int64) error {
	return a.AppStoreRepository.DeleteAppStore(id)
}

// UpdateAppStore 更新应用市场
func (a *AppStore) UpdateAppStore(store *model.AppStore) error {
	return a.AppStoreRepository.UpdateAppStore(store)
}

// FindAppStoreByID 根据ID查询应用市场数据
func (a *AppStore) FindAppStoreByID(id int64) (*model.AppStore, error) {
	return a.AppStoreRepository.FindAppStoreByID(id)
}

// FindAllAppStore 查询全部应用市场数据
func (a *AppStore) FindAllAppStore() ([]model.AppStore, error) {
	return a.AppStoreRepository.FindAll()
}

// AddInstallNum 添加安装数量
func (a *AppStore) AddInstallNum(id int64) error {
	return a.AppStoreRepository.AddInstallNumber(id)
}

// GetInstallNum 查询安装数量
func (a *AppStore) GetInstallNum(id int64) int64 {
	return a.AppStoreRepository.GetInstallNumber(id)
}

// AddViewNum 添加浏览量
func (a *AppStore) AddViewNum(id int64) error {
	return a.AppStoreRepository.AddViewNumber(id)
}

// GetViewNum 查询浏览量
func (a *AppStore) GetViewNum(id int64) int64 {
	return a.AppStoreRepository.GetViewNumber(id)
}

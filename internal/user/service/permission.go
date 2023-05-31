package service

import (
	"k8s.io/client-go/kubernetes"
	"tini-paas/internal/user/model"
	"tini-paas/internal/user/repository"
)

// PermissionService 权限服务接口
type PermissionService interface {
	AddPermission(permission *model.Permission) (int64, error)
	DeletePermission(id int64) error
	UpdatePermission(permission *model.Permission) error
	FindPermissionByID(id int64) (*model.Permission, error)
	FindAllPermission() ([]model.Permission, error)

	// 根据角色ID查找所有权限
	FindAllPermissionByID(id []int64) ([]*model.Permission, error)
}

// NewPermissionService 初始化权限服务
func NewPermissionService(permissionRepository repository.PermissionRepository, client *kubernetes.Clientset) PermissionService {
	return &PermissionDataService{
		PermissionRepository: permissionRepository,
	}
}

// PermissionDataService  权限服务对象
type PermissionDataService struct {
	PermissionRepository repository.PermissionRepository
}

// AddPermission 添加权限
func (p *PermissionDataService) AddPermission(permission *model.Permission) (int64, error) {
	return p.PermissionRepository.CreatePermission(permission)
}

// DeletePermission 删除权限
func (p *PermissionDataService) DeletePermission(id int64) error {
	return p.PermissionRepository.DeletePermission(id)
}

// UpdatePermission 更新权限
func (p *PermissionDataService) UpdatePermission(permission *model.Permission) error {
	return p.PermissionRepository.UpdatePermission(permission)
}

// FindPermissionByID 查询权限信息
func (p *PermissionDataService) FindPermissionByID(id int64) (*model.Permission, error) {
	return p.PermissionRepository.FindPermissionByID(id)
}

// FindAllPermission 查询全部权限信息
func (p *PermissionDataService) FindAllPermission() ([]model.Permission, error) {
	return p.PermissionRepository.FindAll()
}

// FindAllPermissionByID 更加角色ID查询权限
func (p *PermissionDataService) FindAllPermissionByID(id []int64) ([]*model.Permission, error) {
	return p.PermissionRepository.FindAllPermissionByID(id)
}

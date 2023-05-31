package service

import (
	"k8s.io/client-go/kubernetes"
	"tini-paas/internal/user/model"
	"tini-paas/internal/user/repository"
)

// RoleService 角色服务接口
type RoleService interface {
	AddRole(user *model.Role) (int64, error)
	DeleteRole(id int64) error
	UpdateRole(role *model.Role) error
	FindRoleByID(id int64) (*model.Role, error)
	FindAllRole() ([]model.Role, error)

	FindAllRoleByID(id []int64) ([]*model.Role, error)
	AddPermission(role *model.Role, perm []*model.Permission) error
	UpdatePermission(role *model.Role, perm []*model.Permission) error
	DeletePermission(role *model.Role, perm []*model.Permission) error
}

// NewRoleService 初始化角色服务
func NewRoleService(roleRepository repository.RoleRepository, client *kubernetes.Clientset) RoleService {
	return &RoleDataService{
		RoleRepository: roleRepository,
	}
}

// RoleDataService  角色服务操作对象
type RoleDataService struct {
	RoleRepository repository.RoleRepository
}

// AddRole 添加角色
func (r *RoleDataService) AddRole(user *model.Role) (int64, error) {
	return r.RoleRepository.CreateRole(user)
}

// DeleteRole 删除角色
func (r *RoleDataService) DeleteRole(id int64) error {
	return r.RoleRepository.DeleteRole(id)
}

// UpdateRole 更新角色
func (r *RoleDataService) UpdateRole(role *model.Role) error {
	return r.RoleRepository.UpdateRole(role)
}

// FindRoleByID 查询角色信息
func (r *RoleDataService) FindRoleByID(id int64) (*model.Role, error) {
	return r.RoleRepository.FindRoleByID(id)
}

// FindAllRole 查找全部角色信息
func (r *RoleDataService) FindAllRole() ([]model.Role, error) {
	return r.RoleRepository.FindAll()
}

// FindAllRoleByID 根据用户ID查找所有角色
func (r *RoleDataService) FindAllRoleByID(id []int64) ([]*model.Role, error) {
	return r.RoleRepository.FindAllRoleByID(id)
}

// AddPermission 为角色添加权限
func (r *RoleDataService) AddPermission(role *model.Role, perm []*model.Permission) error {
	return r.RoleRepository.AddPermission(role, perm)
}

// UpdatePermission 更新权限
func (r *RoleDataService) UpdatePermission(role *model.Role, perm []*model.Permission) error {
	return r.RoleRepository.UpdatePermission(role, perm)
}

// DeletePermission 删除权限
func (r *RoleDataService) DeletePermission(role *model.Role, perm []*model.Permission) error {
	return r.RoleRepository.DeletePermission(role, perm)
}

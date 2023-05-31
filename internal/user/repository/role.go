package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/user/model"
)

// RoleRepository 角色数据操作接口
type RoleRepository interface {
	InitTable() error
	CreateRole(role *model.Role) (int64, error)
	DeleteRole(id int64) error
	UpdateRole(role *model.Role) error
	FindRoleByID(id int64) (*model.Role, error)
	FindAll() ([]model.Role, error)

	FindAllRoleByID(id []int64) ([]*model.Role, error)
	AddPermission(role *model.Role, perm []*model.Permission) error
	UpdatePermission(role *model.Role, perm []*model.Permission) error
	DeletePermission(role *model.Role, perm []*model.Permission) error
}

// NewRoleRepository 初始化角色
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &Role{
		db: db,
	}
}

// Role 角色数据操作对象
type Role struct {
	db *gorm.DB
}

// InitTable 初始化表
func (r *Role) InitTable() error {
	return r.db.CreateTable(&model.Role{}).Error
}

// CreateRole 创建角色
func (r *Role) CreateRole(role *model.Role) (int64, error) {
	return role.ID, r.db.Create(role).Error
}

// DeleteRole 删除角色
func (r *Role) DeleteRole(id int64) error {
	return r.db.Delete(&model.Role{}).Where("id = ?", id).Error
}

// UpdateRole 更新角色
func (r *Role) UpdateRole(role *model.Role) error {
	return r.db.Model(role).Update(role).Error
}

// FindRoleByID 查询角色信息
func (r *Role) FindRoleByID(id int64) (*model.Role, error) {
	role := &model.Role{}
	return role, r.db.First(role).Where("id = ?", id).Error
}

// FindAll 查找全部角色信息
func (r *Role) FindAll() ([]model.Role, error) {
	var roleAll []model.Role
	return roleAll, r.db.Find(roleAll).Error
}

// FindAllRoleByID 根据ID查找所有角色
func (r *Role) FindAllRoleByID(id []int64) ([]*model.Role, error) {
	var roleAll []*model.Role
	return roleAll, r.db.Find(&roleAll, id).Error
}

// AddPermission 添加角色权限
func (r *Role) AddPermission(role *model.Role, perm []*model.Permission) error {
	return r.db.Model(&role).Association("Permission").Append(perm).Error
}

// UpdatePermission 更新角色权限
func (r *Role) UpdatePermission(role *model.Role, perm []*model.Permission) error {
	return r.db.Model(&role).Association("Permission").Replace(perm).Error
}

// DeletePermission 删除角色权限
func (r *Role) DeletePermission(role *model.Role, perm []*model.Permission) error {
	return r.db.Model(&role).Association("Permission").Delete(perm).Error
}

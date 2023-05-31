package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/user/model"
)

// PermissionRepository  权限数据操作接口
type PermissionRepository interface {
	InitTable() error
	CreatePermission(permission *model.Permission) (int64, error)
	DeletePermission(id int64) error
	UpdatePermission(permission *model.Permission) error
	FindPermissionByID(id int64) (*model.Permission, error)
	FindAll() ([]model.Permission, error)

	FindAllPermissionByID(id []int64) ([]*model.Permission, error)
}

// NewPermissionRepository 初始化权限
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &Permission{
		db: db,
	}
}

// Permission 权限数据操作对象
type Permission struct {
	db *gorm.DB
}

// InitTable 初始化表
func (p *Permission) InitTable() error {
	return p.db.CreateTable(&model.Permission{}).Error
}

// CreatePermission 创建权限
func (p *Permission) CreatePermission(permission *model.Permission) (int64, error) {
	return permission.ID, p.db.Create(permission).Error
}

// DeletePermission 删除权限
func (p *Permission) DeletePermission(id int64) error {
	return p.db.Delete(&model.Permission{}).Where("id = ?", id).Error
}

// UpdatePermission 更新权限
func (p *Permission) UpdatePermission(permission *model.Permission) error {
	return p.db.Model(permission).Update(permission).Error
}

// FindPermissionByID 查询权限信息
func (p *Permission) FindPermissionByID(id int64) (*model.Permission, error) {
	perm := &model.Permission{}
	return perm, p.db.First(perm, id).Error
}

// FindAll 查询所有权限
func (p *Permission) FindAll() ([]model.Permission, error) {
	var permAll []model.Permission
	return permAll, p.db.Find(permAll).Error
}

// FindAllPermissionByID 根据角色id查找所有权限
func (p *Permission) FindAllPermissionByID(id []int64) ([]*model.Permission, error) {
	var permissionAll []*model.Permission
	return permissionAll, p.db.Find(&permissionAll, id).Error
}

package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/user/model"
)

// UserRepository 用户数据操作接口
type UserRepository interface {
	InitTable() error
	CreateUser(user *model.User) (int64, error)
	DeleteUser(id int64) error
	UpdateUser(user *model.User) error
	FindUserByID(id int64) (*model.User, error)
	FindAll() ([]model.User, error)

	AddRole(user *model.User, role []*model.Role) error
	UpdateRole(user *model.User, role []*model.Role) error
	DeleteRole(user *model.User, role []*model.Role) error
	IsRight(action string, id int64) bool
}

// NewUserRepository 初始化用户中心
func NewUserRepository(db *gorm.DB) UserRepository {
	return &User{
		db: db,
	}
}

// User 用户中心数据操作对象
type User struct {
	db *gorm.DB
}

// InitTable 初始化表
func (u *User) InitTable() error {
	return u.db.CreateTable(&model.User{}, &model.Role{}, &model.Permission{}).Error
}

// CreateUser 创建用户
func (u *User) CreateUser(user *model.User) (int64, error) {
	return user.ID, u.db.Create(user).Error
}

// DeleteUser 删除用户
func (u *User) DeleteUser(id int64) error {
	return u.db.Delete(&model.User{}).Where("id = ?", id).Error
}

// UpdateUser 更新用户
func (u *User) UpdateUser(user *model.User) error {
	return u.db.Model(user).Update(user).Error
}

// FindUserByID 查找用户信息
func (u *User) FindUserByID(id int64) (*model.User, error) {
	user := &model.User{}
	return user, u.db.First(user, id).Error
}

// FindAll 查找所有用户信息
func (u *User) FindAll() ([]model.User, error) {
	var users []model.User
	return users, u.db.Find(users).Error
}

// AddRole 添加角色
func (u *User) AddRole(user *model.User, role []*model.Role) error {
	return u.db.Model(&user).Association("Role").Append(role).Error
}

// UpdateRole 更新角色
func (u *User) UpdateRole(user *model.User, role []*model.Role) error {
	return u.db.Model(&user).Association("Role").Replace(role).Error
}

// DeleteRole 删除角色
func (u *User) DeleteRole(user *model.User, role []*model.Role) error {
	return u.db.Model(&user).Association("Role").Delete(role).Error
}

// IsRight 鉴权
func (u *User) IsRight(action string, id int64) bool {
	permission := &model.Permission{}
	sql := "SELECT p.id FROM user u, user_role ur, role r, role_permission rp, permission.proto p WHERE p.permission_action=? AND  p.id = rp.permission_id AND rp.role_id = r.id AND ur.role_id = r.id AND ur.user_id = u.id AND u.id = ?"

	// 原生SQL
	u.db.Raw(sql, action, id).Scan(permission)
	if permission.ID > 0 {
		return true
	}
	return false
}

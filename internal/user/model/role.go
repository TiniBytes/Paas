package model

// Role 角色
type Role struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// RoleName 角色名称
	RoleName string `json:"role_name"`

	// RoleStatus 角色状态
	RoleStatus int32 `json:"role_status"`

	// Permission 角色权限
	Permission []*Permission `gorm:"many2many:role_permission" json:"permission.proto"`
}

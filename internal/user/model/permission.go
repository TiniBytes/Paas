package model

// Permission 权限
type Permission struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// PermissionName 权限名称
	PermissionName string `json:"permission_name"`

	// PermissionDescribe 权限描述
	PermissionDescribe string `json:"permission_describe"`

	// PermissionAction 权限行为
	PermissionAction string `json:"permission_action"`

	// PermissionStatus 权限状态
	PermissionStatus int32 `json:"permission_status"`
}

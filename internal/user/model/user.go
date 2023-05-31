package model

// User 用户
type User struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// UserName 用户名
	UserName string `gorm:"not_null;unique" json:"user_name"`

	// UserEmail 用户邮箱
	UserEmail string `gorm:"not_null;unique" json:"user_email"`

	// IsAdmin 是否管理员
	IsAdmin bool `json:"is_admin"`

	// UserPwd 用户密码
	UserPwd string `json:"user_pwd"`

	// UserStatus 用户状态
	UserStatus int32 `json:"user_status"`

	// Role 角色
	Role []*Role `gorm:"many2many:user_role"`
}

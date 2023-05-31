package model

// MiddleConfig 中间件的初始化信息
type MiddleConfig struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`

	// MiddleID 关联的中间件ID
	MiddleID int64 `json:"middle_id"`

	// MiddleConfigRootUser 可能存在的root用户
	MiddleConfigRootUser string `json:"middle_config_root_user"`

	// MiddleConfigRootPwd 可能存在的root密码
	MiddleConfigRootPwd string `json:"middle_config_root_pwd"`

	// MiddleConfigUser 可能存在的普通用户
	MiddleConfigUser string `json:"middle_config_user"`

	// MiddleConfigPwd 普通用户的密码
	MiddleConfigPwd string `json:"middle_config_user_pwd"`

	// MiddleConfigDataBase 预置数据库名字
	MiddleConfigDataBase string `json:"middle_config_data_base"`
}

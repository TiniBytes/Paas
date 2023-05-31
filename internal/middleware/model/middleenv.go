package model

// MiddleEnv 中间件的环境变量
type MiddleEnv struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`

	// MiddleID 关联的中间件ID
	MiddleID int64 `json:"middle_id"`

	// EnvKey 环境key
	EnvKey string `json:"env_key"`

	// EnvValue 环境value
	EnvValue string `json:"env_value"`
}

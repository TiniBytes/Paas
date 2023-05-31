package model

// PodEnv pod环境变量
type PodEnv struct {
	// ID 主键
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`

	// PodID podApi id
	PodID int64 `json:"pod_id"`

	// EnvKey 环境key
	EnvKey string `json:"env_key"`

	// EnvValue 环境value
	EnvValue string `json:"env_value"`
}

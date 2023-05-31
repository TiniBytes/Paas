package model

// Middleware 动态中间件
type Middleware struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// MiddleName 中间件名称
	MiddleName string `json:"middle_name"`

	// MiddleNamespace 中间件命名空间
	MiddleNamespace string `json:"middle_namespace"`

	// MiddleTypeID 中间件类型
	MiddleTypeID int64 `json:"middle_type_id"`

	// MiddleVersionID 中间件的版本
	MiddleVersionID int64 `json:"middle_version_id"`

	// MiddlePort 中间件的端口
	MiddlePort []MiddlePort `gorm:"ForeignKey:MiddleID" json:"middle_port"`

	// MiddleConfig 默认生成的账号密码
	MiddleConfig MiddleConfig `gorm:"ForeignKey:MiddleID" json:"middle_config"`

	// MiddleEnv 环境变量
	MiddleEnv []MiddleEnv `gorm:"ForeignKey:MiddleID" json:"middle_env"`

	// MiddleCPU 中间件的CPU管控
	MiddleCPU float32 `json:"middle_cpu"`

	// MiddleMemory 中间件的内存管控
	MiddleMemory float32 `json:"middle_memory"`

	// MiddleStorage 中间件存储
	MiddleStorage []MiddleStorage `gorm:"ForeignKey:MiddleID" json:"middle_storage"`

	// MiddleReplicas 中间件副本
	MiddleReplicas int32 `json:"middle_replicas"`
}

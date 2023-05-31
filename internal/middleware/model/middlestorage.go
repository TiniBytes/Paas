package model

// MiddleStorage 中间件存储
type MiddleStorage struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`

	// MiddleID 关联的中间件ID
	MiddleID int64 `json:"middle_id"`

	// MiddleStorageName 存储名称
	MiddleStorageName string `json:"middle_storage_name"`

	// MiddleStorageSize 存储的大小
	MiddleStorageSize float32 `json:"middle_storage_size"`

	// MiddleStoragePath 存储需要挂载的目录
	MiddleStoragePath string `json:"middle_storage_path"`

	// MiddleStorageClass 存储创建的类型
	MiddleStorageClass string `json:"middle_storage_class"`

	// MiddleStorageAccessMode 存储的权限
	MiddleStorageAccessMode string `json:"middle_storage_access_mode"`
}

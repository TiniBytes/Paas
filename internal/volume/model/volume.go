package model

// Volume 存储卷
type Volume struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// VolumeName 存储名称
	VolumeName string `json:"volume_name"`

	// VolumeNamespace 存储所属的命名空间
	VolumeNamespace string `json:"volume_namespace"`

	// VolumeAccessMode 存储的访问模式：RWO, ROX, RWX
	VolumeAccessMode string `json:"volume_access_mode"`

	// VolumeStorageClassName 存储类名称
	VolumeStorageClassName string `json:"volume_storage_class_name"`

	// VolumeRequest 存储请求资源大小
	VolumeRequest float32 `json:"volume_request"`

	// VolumePersistentVolumeMode 存储类型：Block，filesystem
	VolumePersistentVolumeMode string `json:"volume_persistent_volume_mode"`
}

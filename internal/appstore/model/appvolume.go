package model

// AppVolume 应用存储
type AppVolume struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// AppID 关联的应用ID
	AppID int64 `json:"app_id"`

	// AppVolumeID 存储ID
	AppVolumeID int64 `json:"app_volume_id"`
}

package model

// AppImage 应用图片
type AppImage struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// AppID 关联的应用ID
	AppID int64 `json:"app_id"`

	// AppImageSrc 图片地址
	AppImageSrc string `json:"app_image_src"`
}

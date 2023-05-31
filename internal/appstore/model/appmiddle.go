package model

// AppMiddle 云应用中间件模板
type AppMiddle struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// AppID 关联的应用ID
	AppID int64 `json:"app_id"`

	// AppMiddleID 中间件ID
	AppMiddleID int64 `json:"app_middle_id"`
}

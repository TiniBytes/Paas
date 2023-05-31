package model

// AppIsv 应用服务商
type AppIsv struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// AppIsvName 应用服务商名称
	AppIsvName string `json:"app_isv_name"`

	// AppIsvDetail 应用服务商详情
	AppIsvDetail string `json:"app_isv_detail"`
}

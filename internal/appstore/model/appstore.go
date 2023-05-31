package model

// AppStore 云应用市场
type AppStore struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// AppSku 应用唯一标识
	AppSku string `gorm:"unique_index;not_null" json:"app_sku"`

	// AppTitle 应用标题
	AppTitle string `json:"app_title"`

	// AppDetail 应用描述
	AppDetail string `json:"app_detail"`

	// AppPrice 应用价格
	AppPrice float32 `json:"app_price"`

	// AppInstall 安装次数
	AppInstall int64 `json:"app_install"`

	// AppViews 访问次数
	AppViews int64 `json:"app_views"`

	// AppCheck 应用审核
	AppCheck bool `json:"app_check"`

	// AppCategoryID 应用分类
	AppCategoryID int64 `json:"app_category_id"`

	// AppIsvID 服务商ID
	AppIsvID int64 `json:"app_isv_id"`

	// AppImage 应用图片
	AppImage []AppImage `gorm:"ForeignKey:AppID" json:"app_image"`

	// AppPod 应用组合(应用的模板)
	AppPod []AppPod `gorm:"ForeignKey:AppID" json:"app_pod"`

	// AppMiddle 中间件组合
	AppMiddle []AppMiddle `gorm:"ForeignKey:AppID" json:"app_middle"`

	// AppVolume 存储组合
	AppVolume []AppVolume `gorm:"ForeignKey:AppID" json:"app_volume"`

	// AppComment 评论
	AppComment []AppComment `gorm:"ForeignKey:AppID" json:"app_comment"`
}

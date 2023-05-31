package model

// AppCategory 应用分类
type AppCategory struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// CategoryName 分类名称
	CategoryName string `json:"category_name"`
}

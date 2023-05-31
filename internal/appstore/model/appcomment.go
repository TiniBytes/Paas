package model

// AppComment 应用评论
type AppComment struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// AppID 关联的应用ID
	AppID int64 `json:"app_id"`

	// AppCommentTitle 评论标题
	AppCommentTitle string `json:"app_comment_title"`

	// AppCommentDetail 评论详情
	AppCommentDetail string `json:"app_comment_detail"`

	// AppUserID 评论用户
	AppUserID int64 `json:"app_user_id"`
}

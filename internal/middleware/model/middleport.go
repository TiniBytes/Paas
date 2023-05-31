package model

// MiddlePort 中间件端口
type MiddlePort struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// MiddleID 关联中间件的ID
	MiddleID int64 `json:"middle_id"`

	// MiddlePort 中间件开放的端口
	MiddlePort int32 `json:"middle_port"`

	// MiddleProtocol 中间件开放端口协议
	MiddleProtocol string `json:"middle_protocol"`
}

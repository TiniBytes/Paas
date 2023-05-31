package model

// SvcPort 服务端口
type SvcPort struct {
	// ID 服务端口ID
	ID int64 `gorm:"primary_key;not_null;auto-increment"`

	// SvcID 服务端口
	SvcID int64 `json:"service_id"`

	// SvcPort 服务端口
	SvcPort string `json:"service_port"`

	// SvcTargetPort pod中需要映射的port地址
	SvcTargetPort int32 `json:"service_target_port"`

	// SvcNodePort 开启NodePort的模式下进行设置
	SvcNodePort int32 `json:"service_node_port"`

	// SvcPortProtocol 端口协议
	SvcPortProtocol string `json:"service_port_protocol"`
}

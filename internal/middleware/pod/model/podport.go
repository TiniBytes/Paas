package model

// PodPort pod端口
type PodPort struct {
	// ID 主键
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`

	// PodID podApi id
	PodID int64 `json:"pod_id"`

	// ContainerPort 容器端口
	ContainerPort int32 `json:"container_port"`

	// Protocol 协议
	Protocol string `json:"protocol"`
}

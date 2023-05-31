package model

// Svc 服务
type Svc struct {
	// ID 服务ID
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// SvcName 服务名称
	SvcName string `gorm:"unique_index;not_null" json:"service_name"`

	// SvcNamespace 服务名称命名空间
	SvcNamespace string `gorm:"not_null" json:"service_namespace"`

	// SvcPodName 绑定的pod名称
	SvcPodName string `gorm:"not_null" json:"service_pod_name"`

	// SvcType 服务类型 ClusterIP, NodePort, LoadBalancer, ExternalName
	SvcType string `json:"service_type"`

	// SvcExternalName 服务外部名称， ExternalName时候启用该字段
	SvcExternalName string `json:"service_external_name"`

	// SvcTeamID 业务侧团队ID
	SvcTeamID string `json:"service_team_id"`

	// SvcPort 服务上的端口设置
	SvcPort []SvcPort `gorm:"ForeignKey:SvcID" json:"service_port"`
}

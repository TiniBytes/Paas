package model

// RoutePath 关联路径
type RoutePath struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// RouteID 关联RouteID
	RouteID int64 `json:"route_id"`

	// RoutePathName url
	RoutePathName string `json:"route_path_name"`

	// RouteBackendService route绑定service的名称
	RouteBackendService string `json:"route_backend_service"`

	// RouteBackendServicePort route绑定service暴露的端口
	RouteBackendServicePort string `json:"route_backend_service_port"`
}

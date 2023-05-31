package model

// Route 路由
type Route struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// RouteName 路由名称
	RouteName string `json:"route_name"`

	// RouteNamespace 路由命名空间
	RouteNamespace string `json:"route_namespace"`

	// RouteHost 路由域名
	RouteHost string `json:"route_host"`

	// RoutePath 关联路径
	RoutePath []RoutePath `gorm:"ForeignKey:RouteID" json:"route_path"`
}

package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/route/model"
	"tini-paas/pkg/common"
)

// RouteRepository route操作
type RouteRepository interface {
	// InitTable 初始化表
	InitTable() error

	// CreateRoute 创建Route
	CreateRoute(*model.Route) (int64, error)

	// DeleteRouteByID 删除Route
	DeleteRouteByID(int64) error

	// UpdateRoute 更新Route
	UpdateRoute(*model.Route) error

	// FindRouteByID 查找Route
	FindRouteByID(int64) (*model.Route, error)

	// FindAll 查找所有Route
	FindAll() ([]model.Route, error)
}

// NewRouteRepository 创建Route对象
func NewRouteRepository(db *gorm.DB) RouteRepository {
	return &Route{
		db: db,
	}
}

// Route pod repository
type Route struct {
	db *gorm.DB
}

// InitTable 初始化表
func (r *Route) InitTable() error {
	return r.db.CreateTable(&model.Route{}, &model.RoutePath{}).Error
}

// CreateRoute 创建Route
func (r *Route) CreateRoute(route *model.Route) (int64, error) {
	return route.ID, r.db.Create(route).Error
}

// DeleteRouteByID 删除Route
func (r *Route) DeleteRouteByID(i int64) error {
	// 开始事务
	tx := r.db.Begin()
	// 遇到问题回滚
	defer func() {
		if re := recover(); re != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		common.Error(tx.Error)
		return tx.Error
	}

	// 开始执行删除
	err := r.db.Delete(&model.Route{}).Where("id = ", i).Error
	if err != nil {
		tx.Rollback()
		common.Error(tx.Error)
		return err
	}

	err = r.db.Delete(&model.RoutePath{}).Where("route_id", i).Error
	if err != nil {
		tx.Rollback()
		common.Error(tx.Error)
		return err
	}

	return tx.Commit().Error
}

// UpdateRoute 更新Route
func (r *Route) UpdateRoute(route *model.Route) error {
	return r.db.Model(route).Update(route).Error
}

// FindRouteByID 查找Route
func (r *Route) FindRouteByID(i int64) (*model.Route, error) {
	route := &model.Route{}
	return route, r.db.Preload("RoutePath").First(route, i).Error
}

// FindAll 查找所有Route
func (r *Route) FindAll() ([]model.Route, error) {
	var routeAll []model.Route
	return routeAll, r.db.Preload("RoutePath").Find(&routeAll).Error
}

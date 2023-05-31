package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/middleware/model"
)

// MiddlewareRepository 中间件数据库操作接口
type MiddlewareRepository interface {
	InitTable() error
	CreateMiddleware(*model.Middleware) (int64, error)
	DeleteMiddleware(int64) error
	UpdateMiddleware(*model.Middleware) error
	FindMiddlewareByID(int64) (*model.Middleware, error)
	FindAll() ([]model.Middleware, error)

	// FindAllByTypeID 根据类型查找中间件
	FindAllByTypeID(int64) ([]model.Middleware, error)
}

// NewMiddlewareRepository 初始化中间件
func NewMiddlewareRepository(db *gorm.DB) MiddlewareRepository {
	return &Middleware{
		db: db,
	}
}

// Middleware 中间件数据库操作对象
type Middleware struct {
	db *gorm.DB
}

func (m *Middleware) InitTable() error {
	return m.db.CreateTable(&model.Middleware{}, &model.MiddleConfig{}, &model.MiddlePort{}, &model.MiddleEnv{}, &model.MiddleStorage{}).Error
}

func (m *Middleware) CreateMiddleware(middleware *model.Middleware) (int64, error) {
	return middleware.ID, m.db.Create(middleware).Error
}

func (m *Middleware) DeleteMiddleware(i int64) error {
	// 开启事务
	tx := m.db.Begin()
	defer func() {
		// 遇到问题回滚
		if r := recover(); r != nil {
			tx.Callback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除中间件
	err := m.db.Delete(&model.Middleware{}).Where("id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除中间件配置
	err = m.db.Delete(&model.MiddleConfig{}).Where("middle_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除端口
	err = m.db.Delete(&model.MiddlePort{}).Where("middle_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除环境变量
	err = m.db.Delete(&model.MiddleEnv{}).Where("middle_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除中间件存储
	err = m.db.Delete(&model.MiddleStorage{}).Where("middle_id = ?", i).Error
	if err != nil {
		tx.Callback()
		return err
	}

	return tx.Commit().Error
}

func (m *Middleware) UpdateMiddleware(middleware *model.Middleware) error {
	return m.db.Model(middleware).Update(middleware).Error
}

func (m *Middleware) FindMiddlewareByID(i int64) (*model.Middleware, error) {
	middleware := &model.Middleware{}
	return middleware, m.db.First(&middleware, i).Error
}

func (m *Middleware) FindAll() ([]model.Middleware, error) {
	var middleAll []model.Middleware
	return middleAll, m.db.Find(&middleAll).Error
}

func (m *Middleware) FindAllByTypeID(i int64) ([]model.Middleware, error) {
	var middleAll []model.Middleware
	return middleAll, m.db.Find(&middleAll).Where("middle_type_id = ?", i).Error
}

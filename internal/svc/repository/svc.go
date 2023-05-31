package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/svc/model"
	"tini-paas/pkg/common"
)

// SvcRepository Service接口
type SvcRepository interface {
	// InitTable 初始化表
	InitTable() error

	// CreateSvc 创建一条service数据
	CreateSvc(*model.Svc) (int64, error)

	// DeleteSvcByID 根据ID删除一条service数据
	DeleteSvcByID(int64) error

	// UpdateSvc 更新数据
	UpdateSvc(*model.Svc) error

	// FindSvcByID 根据ID查找数据
	FindSvcByID(int64) (*model.Svc, error)

	// FindAll 查找所有service数据
	FindAll() ([]model.Svc, error)
}

// NewSvcRepository 初始化ServiceRepository
func NewSvcRepository(db *gorm.DB) SvcRepository {
	return &Svc{
		db: db,
	}
}

// Svc 服务repository
type Svc struct {
	db *gorm.DB
}

// InitTable 初始化表
func (s *Svc) InitTable() error {
	// 初始化表结构，创建Service表和ServicePort表
	return s.db.CreateTable(&model.Svc{}, &model.SvcPort{}).Error
}

// CreateSvc 创建一条service数据
func (s *Svc) CreateSvc(service *model.Svc) (int64, error) {
	return service.ID, s.db.Create(service).Error
}

// DeleteSvcByID 根据ID删除一条service数据
func (s *Svc) DeleteSvcByID(i int64) error {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			// 遇到错误回滚
			tx.Callback()
		}
	}()
	if tx.Error != nil {
		common.Error(tx.Error)
		return tx.Error
	}

	// 删除service
	err := s.db.Delete(&model.Svc{}).Where("id = ?", i).Error
	if err != nil {
		tx.Callback()
		common.Error(err)
		return err
	}

	// 删除相关的port
	err = s.db.Delete(&model.SvcPort{}).Where("service_id = ?", i).Error
	if err != nil {
		tx.Callback()
		common.Error(err)
		return err
	}
	return tx.Commit().Error
}

// UpdateSvc 更新数据
func (s *Svc) UpdateSvc(service *model.Svc) error {
	return s.db.Model(service).Update(service).Error
}

// FindSvcByID 根据ID查找数据
func (s *Svc) FindSvcByID(i int64) (*model.Svc, error) {
	// 将数据写入到service中返回
	service := &model.Svc{}
	return service, s.db.Preload("SvcPort").First(service, i).Error
}

// FindAll 查找所有service数据
func (s *Svc) FindAll() ([]model.Svc, error) {
	var serviceAll []model.Svc
	return serviceAll, s.db.Find(&serviceAll).Error
}

package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/middleware/model"
)

// MiddleTypeRepository 中间件类型接口
type MiddleTypeRepository interface {
	InitTable() error
	CreateMiddleType(*model.MiddleType) (int64, error)
	DeleteMiddleType(int64) error
	UpdateMiddleType(*model.MiddleType) error
	FindMiddleTypeByID(int64) (*model.MiddleType, error)
	FindAll() ([]model.MiddleType, error)

	// FindVersionByID 查询版本信息
	FindVersionByID(int64) (*model.MiddleVersion, error)

	// FindAllVersionByTypeID 根据类型查询全部版本信息
	FindAllVersionByTypeID(int64) ([]model.MiddleVersion, error)
}

// NewMiddleTypeRepository 初始化中间件类型
func NewMiddleTypeRepository(db *gorm.DB) MiddleTypeRepository {
	return &MiddleType{
		db: db,
	}
}

// MiddleType 中间件类型接口
type MiddleType struct {
	db *gorm.DB
}

func (m *MiddleType) InitTable() error {
	return m.db.CreateTable(&model.MiddleType{}, &model.MiddleVersion{}).Error
}

func (m *MiddleType) CreateMiddleType(middleType *model.MiddleType) (int64, error) {
	return middleType.ID, m.db.Create(middleType).Error
}

func (m *MiddleType) DeleteMiddleType(i int64) error {
	// 开启事务
	tx := m.db.Begin()
	defer func() {
		// 遇到问题回滚
		if r := recover(); r != nil {
			tx.Callback()
		}
	}()
	if tx.Error != nil {
		tx.Callback()
		return nil
	}

	// 开始删除
	err := m.db.Delete(&model.MiddleType{}).Where("id = ?").Error
	if err != nil {
		tx.Callback()
		return err
	}

	// 删除版本信息
	err = m.db.Delete(&model.MiddleVersion{}).Where("middle_type_id").Error
	if err != nil {
		tx.Callback()
		return err
	}

	return tx.Commit().Error
}

func (m *MiddleType) UpdateMiddleType(middleType *model.MiddleType) error {
	return m.db.Model(middleType).Update(middleType).Error
}

func (m *MiddleType) FindMiddleTypeByID(i int64) (*model.MiddleType, error) {
	middleType := &model.MiddleType{}
	return middleType, m.db.Preload("MiddleVersion").First(middleType, i).Error
}

func (m *MiddleType) FindAll() ([]model.MiddleType, error) {
	var middleTypeAll []model.MiddleType
	return middleTypeAll, m.db.Find(&middleTypeAll).Error
}

func (m *MiddleType) FindVersionByID(i int64) (*model.MiddleVersion, error) {
	middleVersion := &model.MiddleVersion{}
	return middleVersion, m.db.First(middleVersion, i).Error
}

func (m *MiddleType) FindAllVersionByTypeID(i int64) ([]model.MiddleVersion, error) {
	var middleVersionAll []model.MiddleVersion
	return middleVersionAll, m.db.Find(&middleVersionAll).Where("middle_type_id = ?", i).Error
}

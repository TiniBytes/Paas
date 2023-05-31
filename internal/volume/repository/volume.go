package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/volume/model"
)

// VolumeRepository 存储卷数据库操作接口
type VolumeRepository interface {
	InitTable() error
	CreateVolume(*model.Volume) (int64, error)
	DeleteVolume(int64) error
	UpdateVolume(*model.Volume) error
	FindVolumeByID(int64) (*model.Volume, error)
	FindAll() ([]model.Volume, error)
}

// NewVolumeRepository 初始化数据操作对象
func NewVolumeRepository(db *gorm.DB) VolumeRepository {
	return &Volume{
		db: db,
	}
}

// Volume 数据库对象
type Volume struct {
	db *gorm.DB
}

func (v *Volume) InitTable() error {
	return v.db.CreateTable(&model.Volume{}).Error
}

func (v *Volume) CreateVolume(volume *model.Volume) (int64, error) {
	return volume.ID, v.db.Create(volume).Error
}

func (v *Volume) DeleteVolume(i int64) error {
	return v.db.Delete(&model.Volume{}).Where("id = ?", i).Error
}

func (v *Volume) UpdateVolume(volume *model.Volume) error {
	return v.db.Model(volume).Update(volume).Error
}

func (v *Volume) FindVolumeByID(i int64) (*model.Volume, error) {
	volume := &model.Volume{}
	return volume, v.db.First(volume, i).Error
}

func (v *Volume) FindAll() ([]model.Volume, error) {
	var volumeAll []model.Volume
	return volumeAll, v.db.Find(&volumeAll).Error
}

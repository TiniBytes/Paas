package repository

import (
	"github.com/jinzhu/gorm"
	"tini-paas/internal/pod/model"
)

// PodRepository pod操作
type PodRepository interface {
	// InitTable 初始化表
	InitTable() error

	// FindPodByID 查找pod
	FindPodByID(int64) (*model.Pod, error)

	// CreatePod 创建pod
	CreatePod(*model.Pod) (int64, error)

	// DeletePodByID 删除pod
	DeletePodByID(int64) error

	// UpdatePod 修改pod
	UpdatePod(*model.Pod) error

	// FindAll 查找所有pod
	FindAll() ([]model.Pod, error)
}

// Pod podApi repository
type Pod struct {
	db *gorm.DB
}

// NewPodRepository 创建pod
func NewPodRepository(db *gorm.DB) PodRepository {
	return &Pod{
		db: db,
	}
}

// InitTable 初始化表
func (p *Pod) InitTable() error {
	// 创建三个表
	return p.db.CreateTable(&model.Pod{}, &model.PodEnv{}, &model.PodPort{}).Error
	//return p.db.CreateTable(&model.Pod{}, &model.PodPort{}, &model.PodEnv{}).Error
}

// FindPodByID 查找pod
func (p *Pod) FindPodByID(i int64) (*model.Pod, error) {
	pod := &model.Pod{}
	return pod, p.db.Preload("PodEnv").Preload("PodPort").First(pod, i).Error
}

// CreatePod 创建pod
func (p *Pod) CreatePod(pod *model.Pod) (int64, error) {
	return pod.ID, p.db.Create(pod).Error
}

// DeletePodByID 删除pod
func (p *Pod) DeletePodByID(i int64) error {
	// 开启事务
	tx := p.db.Begin()
	defer func() {
		// 遇到问题回滚
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除pod信息
	err := p.db.Delete(&model.Pod{}).Where("id = ?", i).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除podPort信息
	err = p.db.Delete(&model.PodPort{}).Where("pod_id = ?", i).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除podEnv信息
	err = p.db.Delete(&model.PodEnv{}).Delete("pod_id = ?", i).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 提交
	return tx.Commit().Error
}

// UpdatePod 更新pod
func (p *Pod) UpdatePod(pod *model.Pod) error {
	return p.db.Model(pod).Update(pod).Error
}

// FindAll 获取结果集合
func (p *Pod) FindAll() ([]model.Pod, error) {
	var podAll []model.Pod
	return podAll, p.db.Find(&podAll).Error
}

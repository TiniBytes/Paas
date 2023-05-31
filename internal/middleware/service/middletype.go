package service

import (
	"tini-paas/internal/middleware/model"
	"tini-paas/internal/middleware/repository"
)

// MiddleTypeService 中间件类型服务接口
type MiddleTypeService interface {
	AddMiddleType(*model.MiddleType) (int64, error)
	DeleteMiddleType(int64) error
	UpdateMiddleType(*model.MiddleType) error
	FindMiddleTypeByID(int64) (*model.MiddleType, error)
	FindAllMiddleType() ([]model.MiddleType, error)

	// FindImageVersionByID 根据ID查询镜像地址
	FindImageVersionByID(int64) (string, error)
	FindVersionByID(int64) (*model.MiddleVersion, error)
	FindAllVersionByTypeID(int64) ([]model.MiddleVersion, error)
}

// NewMiddleTypeService 初始化中间件类型服务
func NewMiddleTypeService(repository repository.MiddleTypeRepository) MiddleTypeService {
	return &MiddleTypeDataService{
		MiddleTypeRepository: repository,
	}
}

// MiddleTypeDataService 中间件类型对象
type MiddleTypeDataService struct {
	MiddleTypeRepository repository.MiddleTypeRepository
}

func (m *MiddleTypeDataService) AddMiddleType(middleType *model.MiddleType) (int64, error) {
	return m.MiddleTypeRepository.CreateMiddleType(middleType)
}

func (m *MiddleTypeDataService) DeleteMiddleType(i int64) error {
	return m.MiddleTypeRepository.DeleteMiddleType(i)
}

func (m *MiddleTypeDataService) UpdateMiddleType(middleType *model.MiddleType) error {
	return m.MiddleTypeRepository.UpdateMiddleType(middleType)
}

func (m *MiddleTypeDataService) FindMiddleTypeByID(i int64) (*model.MiddleType, error) {
	return m.MiddleTypeRepository.FindMiddleTypeByID(i)
}

func (m *MiddleTypeDataService) FindAllMiddleType() ([]model.MiddleType, error) {
	return m.MiddleTypeRepository.FindAll()
}

func (m *MiddleTypeDataService) FindImageVersionByID(i int64) (string, error) {
	version, err := m.MiddleTypeRepository.FindVersionByID(i)
	if err != nil {
		return "", err
	}

	// 返回需要的镜像地址
	return version.MiddleDockerImage + ":" + version.MiddleVersion, nil
}

func (m *MiddleTypeDataService) FindVersionByID(i int64) (*model.MiddleVersion, error) {
	return m.MiddleTypeRepository.FindVersionByID(i)
}

func (m *MiddleTypeDataService) FindAllVersionByTypeID(i int64) ([]model.MiddleVersion, error) {
	return m.MiddleTypeRepository.FindAllVersionByTypeID(i)
}

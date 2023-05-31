package service

import (
	"k8s.io/client-go/kubernetes"
	"tini-paas/internal/user/model"
	"tini-paas/internal/user/repository"
)

// UserService 用户服务接口
type UserService interface {
	AddUser(user *model.User) (int64, error)
	DeleteUser(id int64) error
	UpdateUser(user *model.User) error
	FindUserByID(id int64) (*model.User, error)
	FindAllUser() ([]model.User, error)

	AddRole(user *model.User, role []*model.Role) error
	DeleteRole(user *model.User, role []*model.Role) error
	UpdateRole(user *model.User, role []*model.Role) error
	IsRight(action string, id int64) bool
}

// NewUserService 初始化用户服务
func NewUserService(userRepository repository.UserRepository, client *kubernetes.Clientset) UserService {
	return &UserDataService{
		UserRepository: userRepository,
	}
}

// UserDataService 用户服务操作对象
type UserDataService struct {
	UserRepository repository.UserRepository
}

// AddUser 添加用户
func (u *UserDataService) AddUser(user *model.User) (int64, error) {
	return u.UserRepository.CreateUser(user)
}

// DeleteUser 删除用户
func (u *UserDataService) DeleteUser(id int64) error {
	return u.DeleteUser(id)
}

// UpdateUser 更新用户
func (u *UserDataService) UpdateUser(user *model.User) error {
	return u.UserRepository.UpdateUser(user)
}

// FindUserByID 查询用户信息
func (u *UserDataService) FindUserByID(id int64) (*model.User, error) {
	return u.UserRepository.FindUserByID(id)
}

// FindAllUser 查询所有用户
func (u *UserDataService) FindAllUser() ([]model.User, error) {
	return u.UserRepository.FindAll()
}

// AddRole 为用户分配角色
func (u *UserDataService) AddRole(user *model.User, role []*model.Role) error {
	return u.UserRepository.AddRole(user, role)
}

// DeleteRole 删除角色
func (u *UserDataService) DeleteRole(user *model.User, role []*model.Role) error {
	return u.UserRepository.DeleteRole(user, role)
}

// UpdateRole 为用户更新角色
func (u *UserDataService) UpdateRole(user *model.User, role []*model.Role) error {
	return u.UserRepository.UpdateRole(user, role)
}

// IsRight 鉴权
func (u *UserDataService) IsRight(action string, id int64) bool {
	return u.UserRepository.IsRight(action, id)
}

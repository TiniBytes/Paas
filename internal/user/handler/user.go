package handler

import (
	"context"
	"strconv"
	"tini-paas/internal/user/model"
	"tini-paas/internal/user/proto/user"
	"tini-paas/internal/user/service"
	"tini-paas/pkg/common"
)

// UserHandler 用户（API接口实现）
type UserHandler struct {
	UserDataService service.UserService
	RoleDataService service.RoleService
}

// AddUser 添加用户
func (u *UserHandler) AddUser(ctx context.Context, info *user.UserInfo, response *user.Response) error {
	userModel := &model.User{}

	// 数据转化
	err := common.SwapTo(info, userModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := u.UserDataService.AddUser(userModel)
	if err != nil {
		common.Error(err)
		return err
	}
	response.Msg = strconv.FormatInt(rsp, 10)
	return nil
}

// DeleteUser 删除用户
func (u *UserHandler) DeleteUser(ctx context.Context, id *user.UserID, response *user.Response) error {
	return u.UserDataService.DeleteUser(id.Id)
}

// UpdateUser 更新用户
func (u *UserHandler) UpdateUser(ctx context.Context, info *user.UserInfo, response *user.Response) error {
	// 先查询之前是否存在
	userModel, err := u.UserDataService.FindUserByID(info.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将info映射打userModel
	err = common.SwapTo(info, userModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	return u.UserDataService.UpdateUser(userModel)
}

// FindUserByID 查询用户
func (u *UserHandler) FindUserByID(ctx context.Context, id *user.UserID, info *user.UserInfo) error {
	userModel, err := u.UserDataService.FindUserByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	return common.SwapTo(userModel, info)
}

// FindAllUser 查询全部用户
func (u *UserHandler) FindAllUser(ctx context.Context, findAll *user.FindAll, all *user.AllUser) error {
	allUser, err := u.UserDataService.FindAllUser()
	if err != nil {
		common.Error(err)
		return err
	}

	for _, m := range allUser {
		userInfo := &user.UserInfo{}
		err = common.SwapTo(m, userInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		all.UserInfo = append(all.UserInfo, userInfo)
	}
	return nil
}

// AddRole 添加角色
func (u *UserHandler) AddRole(ctx context.Context, role *user.UserRole, response *user.Response) error {
	// 查找用户和对应的角色
	userRole, roles, err := u.getUserRole(role)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 调用后端服务添加角色
	err = u.UserDataService.AddRole(userRole, roles)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

// UpdateRole 更新角色
func (u *UserHandler) UpdateRole(ctx context.Context, role *user.UserRole, response *user.Response) error {
	// 查询用户和对应的角色
	userRole, roles, err := u.getUserRole(role)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 调用后端服务更新角色
	err = u.UserDataService.UpdateRole(userRole, roles)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

// DeleteRole 删除角色
func (u *UserHandler) DeleteRole(ctx context.Context, role *user.UserRole, response *user.Response) error {
	// 查询用户和对应的角色
	userRole, roles, err := u.getUserRole(role)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 调用后端服务删除角色
	err = u.UserDataService.DeleteRole(userRole, roles)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

// IsRight 鉴权
func (u *UserHandler) IsRight(ctx context.Context, userRight *user.UserRight, right *user.Right) error {
	right.Access = u.UserDataService.IsRight(userRight.Action, userRight.UserId)
	return nil
}

// getUserRole 获取用户角色信息
func (u *UserHandler) getUserRole(userRole *user.UserRole) (*model.User, []*model.Role, error) {
	user := &model.User{}
	var role []*model.Role

	// 获取user信息
	user, err := u.UserDataService.FindUserByID(userRole.UserId)
	if err != nil {
		common.Error(err)
		return user, role, err
	}

	// 根据用户信息获取角色信息
	role, err = u.RoleDataService.FindAllRoleByID(userRole.RoleId)
	if err != nil {
		common.Error(err)
		return user, role, err
	}
	return user, role, nil
}

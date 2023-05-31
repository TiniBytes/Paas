package handler

import (
	"context"
	"strconv"
	"tini-paas/internal/user/model"
	"tini-paas/internal/user/proto/role"
	"tini-paas/internal/user/service"
	"tini-paas/pkg/common"
)

// RoleHandler  角色（API接口实现）
type RoleHandler struct {
	RoleDataService       service.RoleService
	PermissionDataService service.PermissionService
}

// AddRole 添加角色
func (r *RoleHandler) AddRole(ctx context.Context, info *role.RoleInfo, response *role.Response) error {
	roleModel := &model.Role{}
	err := common.SwapTo(info, roleModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := r.RoleDataService.AddRole(roleModel)
	if err != nil {
		common.Error(err)
		return err
	}

	response.Msg = strconv.FormatInt(rsp, 10)
	return nil
}

// DeleteRole 删除角色
func (r *RoleHandler) DeleteRole(ctx context.Context, id *role.RoleID, response *role.Response) error {
	return r.RoleDataService.DeleteRole(id.Id)
}

// UpdateRole 更新角色
func (r *RoleHandler) UpdateRole(ctx context.Context, info *role.RoleInfo, response *role.Response) error {
	// 先查询相关数据
	roleModel, err := r.RoleDataService.FindRoleByID(info.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将info映射到roleModel
	err = common.SwapTo(info, roleModel)
	if err != nil {
		common.Error(err)
		return err
	}

	return r.RoleDataService.UpdateRole(roleModel)
}

// FindRoleByID 查询角色
func (r *RoleHandler) FindRoleByID(ctx context.Context, id *role.RoleID, info *role.RoleInfo) error {
	roleModel, err := r.RoleDataService.FindRoleByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将roleModel数据映射到info
	return common.SwapTo(roleModel, info)
}

// FindAllRole 查询全部角色信息
func (r *RoleHandler) FindAllRole(ctx context.Context, findAll *role.FindAll, all *role.AllRole) error {
	allRole, err := r.RoleDataService.FindAllRole()
	if err != nil {
		common.Error(err)
		return err
	}

	for _, v := range allRole {
		roleInfo := &role.RoleInfo{}
		err = common.SwapTo(v, roleInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		all.RoleInfo = append(all.RoleInfo, roleInfo)
	}
	return nil
}

// AddPermission 为角色添加权限
func (r *RoleHandler) AddPermission(ctx context.Context, rolePermission *role.RolePermission, response *role.Response) error {
	role, permission, err := r.getRolePermission(rolePermission)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	err = r.RoleDataService.AddPermission(role, permission)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}
	return nil
}

// DeletePermission 删除权限
func (r *RoleHandler) DeletePermission(ctx context.Context, rolePermission *role.RolePermission, response *role.Response) error {
	role, permission, err := r.getRolePermission(rolePermission)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	err = r.RoleDataService.DeletePermission(role, permission)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}
	return nil
}

// UpdatePermission 更新权限
func (r *RoleHandler) UpdatePermission(ctx context.Context, rolePermission *role.RolePermission, response *role.Response) error {
	role, permission, err := r.getRolePermission(rolePermission)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	err = r.RoleDataService.UpdatePermission(role, permission)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}
	return nil
}

// getRolePermission 获取角色权限信息
func (r *RoleHandler) getRolePermission(rolePermission *role.RolePermission) (role *model.Role, permission []*model.Permission, err error) {
	// 先查询角色ID
	role, err = r.RoleDataService.FindRoleByID(rolePermission.RoleId)
	if err != nil {
		common.Error(err)
		return
	}

	permission, err = r.PermissionDataService.FindAllPermissionByID(rolePermission.PermissionId)
	if err != nil {
		common.Error(err)
		return
	}
	return
}

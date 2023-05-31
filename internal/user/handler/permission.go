package handler

import (
	"context"
	"strconv"
	"tini-paas/internal/user/model"
	"tini-paas/internal/user/proto/permission"
	"tini-paas/internal/user/service"
	"tini-paas/pkg/common"
)

// PermissionHandler 权限（API接口实现）
type PermissionHandler struct {
	PermissionDataService service.PermissionService
}

// AddPermission 添加权限
func (p *PermissionHandler) AddPermission(ctx context.Context, info *permission.PermissionInfo, response *permission.Response) error {
	permissionModel := &model.Permission{}
	err := common.SwapTo(info, permissionModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := p.PermissionDataService.AddPermission(permissionModel)
	if err != nil {
		common.Error(err)
		return err
	}

	response.Msg = strconv.FormatInt(rsp, 10)
	return nil
}

// DeletePermission 删除权限
func (p *PermissionHandler) DeletePermission(ctx context.Context, id *permission.PermissionID, response *permission.Response) error {
	return p.PermissionDataService.DeletePermission(id.Id)
}

// UpdatePermission 更新权限
func (p *PermissionHandler) UpdatePermission(ctx context.Context, info *permission.PermissionInfo, response *permission.Response) error {
	// 先查询权限ID
	permissionModel, err := p.PermissionDataService.FindPermissionByID(info.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	return p.PermissionDataService.UpdatePermission(permissionModel)
}

// FindPermissionByID 查询权限
func (p *PermissionHandler) FindPermissionByID(ctx context.Context, id *permission.PermissionID, info *permission.PermissionInfo) error {
	permissionModel, err := p.PermissionDataService.FindPermissionByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	return common.SwapTo(permissionModel, info)
}

// FindAllPermission 查询所有权限
func (p *PermissionHandler) FindAllPermission(ctx context.Context, findAll *permission.FindAll, all *permission.AllPermission) error {
	allPermission, err := p.PermissionDataService.FindAllPermission()
	if err != nil {
		common.Error(err)
		return err
	}

	for _, v := range allPermission {
		permissionInfo := &permission.PermissionInfo{}
		err = common.SwapTo(v, permissionInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		all.PermissionInfo = append(all.PermissionInfo, permissionInfo)
	}
	return nil
}

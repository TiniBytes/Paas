package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"tini-paas/api/userapi/proto/userApi"
	"tini-paas/internal/user/proto/role"
	"tini-paas/internal/user/proto/user"
	"tini-paas/pkg/common"
)

// UserApi 中间件API处理（对API后端接口方法的实现）
type UserApi struct {
	UserService user.UserService
	RoleService role.RoleService
}

// AddUser 添加用户
func (u *UserApi) AddUser(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	addUser := &user.UserInfo{}

	// 设置user_name
	addUser.UserName, _ = u.getPost(request, "user_name")
	// 设置user_email
	addUser.UserEmail, _ = u.getPost(request, "user_email")
	// 设置is_admin
	admin, _ := u.getPost(request, "is_admin")
	addUser.IsAdmin, _ = strconv.ParseBool(admin)
	// 设置user_pwd
	addUser.UserPwd, _ = u.getPost(request, "user_pwd")

	// 调用后端服务添加
	rsp, err := u.UserService.AddUser(ctx, addUser)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// DeleteUser 删除用户
func (u *UserApi) DeleteUser(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	// 获取用户ID
	userID, err := u.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务删除
	rsp, err := u.UserService.DeleteUser(ctx, &user.UserID{
		Id: userID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// UpdateUser 更新用户
func (u *UserApi) UpdateUser(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	// 获取用户ID
	userID, err := u.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取user_info
	info, err := u.UserService.FindUserByID(ctx, &user.UserID{
		Id: userID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 更新info信息
	info.UserName, _ = u.getPost(request, "user_name")
	info.UserEmail, _ = u.getPost(request, "user_email")
	admin, _ := u.getPost(request, "is_admin")
	info.IsAdmin, _ = strconv.ParseBool(admin)
	info.UserPwd, _ = u.getPost(request, "user_pwd")

	// 调用后端服务执行更新
	rsp, err := u.UserService.UpdateUser(ctx, info)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// FindUserByID 查询用户信息
func (u *UserApi) FindUserByID(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	userID, err := u.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行
	rsp, err := u.UserService.DeleteUser(ctx, &user.UserID{
		Id: userID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// Call 查询全部用户信息（默认方法）
func (u *UserApi) Call(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	allUser, err := u.UserService.FindAllUser(ctx, &user.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(allUser)
	response.Body = string(bytes)
	return nil
}

// AddRole 为用户添加角色
func (u *UserApi) AddRole(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	// 用户ID
	userID, err := u.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 角色ID
	if _, ok := request.Post["role_id"]; !ok {
		common.Error(err)
		return err
	}

	var roleID []int64
	for _, value := range request.Post["role_id"].Values {
		roleID = append(roleID, u.getStringInt64(value))
	}

	// 调用后端执行
	rsp, err := u.UserService.AddRole(ctx, &user.UserRole{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// DeleteRole 删除角色
func (u *UserApi) DeleteRole(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	// 用户ID
	userID, err := u.getID(request)
	if err != nil {
		common.Error(err)
		return err
	}

	// 角色ID
	if _, ok := request.Post["role_id"]; !ok {
		common.Error(err)
		return err
	}

	var roleID []int64
	for _, value := range request.Post["role_id"].Values {
		roleID = append(roleID, u.getStringInt64(value))
	}

	// 调用后端执行
	rsp, err := u.UserService.DeleteRole(ctx, &user.UserRole{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// UpdateRole 更新角色
func (u *UserApi) UpdateRole(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	// 用户ID
	userIDString, err := u.getPost(request, "user_id")
	if err != nil {
		common.Error(err)
		return err
	}
	userID := u.getStringInt64(userIDString)

	// 角色ID
	if _, ok := request.Post["role_id"]; !ok {
		common.Error(err)
		return err
	}

	var roleID []int64
	for _, value := range request.Post["role_id"].Values {
		roleID = append(roleID, u.getStringInt64(value))
	}

	// 调用后端执行
	rsp, err := u.UserService.UpdateRole(ctx, &user.UserRole{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// IsRight 鉴权
func (u *UserApi) IsRight(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	// 获取userID
	if _, ok := request.Get["user_id"]; !ok {
		return errors.New("参数异常")
	}

	idString := request.Get["user_id"].Values[0]
	userID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取action
	if _, ok := request.Get["user_action"]; !ok {
		return errors.New("参数异常")
	}

	userAction := request.Get["user_action"].Values[0]

	// 调用后端服务执行
	right, err := u.UserService.IsRight(ctx, &user.UserRight{
		UserId: userID,
		Action: userAction,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(right)
	response.Body = string(bytes)
	return nil
}

// AddPermission 添加权限
func (u *UserApi) AddPermission(ctx context.Context, request *userApi.Request, response *userApi.Response) error { // 用户ID
	// 角色ID
	roleIDString, err := u.getPost(request, "role_id")
	if err != nil {
		common.Error(err)
		return err
	}
	roleID := u.getStringInt64(roleIDString)

	// 权限ID
	if _, ok := request.Post["permission_id"]; !ok {
		common.Error(err)
		return err
	}

	var permissionID []int64
	for _, value := range request.Post["permission_id"].Values {
		permissionID = append(permissionID, u.getStringInt64(value))
	}

	// 调用后端执行
	rsp, err := u.RoleService.AddPermission(ctx, &role.RolePermission{
		RoleId:       roleID,
		PermissionId: permissionID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// DeletePermission 删除权限
func (u *UserApi) DeletePermission(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	// 角色ID
	roleIDString, err := u.getPost(request, "role_id")
	if err != nil {
		common.Error(err)
		return err
	}
	roleID := u.getStringInt64(roleIDString)

	// 权限ID
	if _, ok := request.Post["permission_id"]; !ok {
		common.Error(err)
		return err
	}

	var permissionID []int64
	for _, value := range request.Post["permission_id"].Values {
		permissionID = append(permissionID, u.getStringInt64(value))
	}

	// 调用后端执行
	rsp, err := u.RoleService.DeletePermission(ctx, &role.RolePermission{
		RoleId:       roleID,
		PermissionId: permissionID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// UpdatePermission 更新权限
func (u *UserApi) UpdatePermission(ctx context.Context, request *userApi.Request, response *userApi.Response) error {
	// 角色ID
	roleIDString, err := u.getPost(request, "role_id")
	if err != nil {
		common.Error(err)
		return err
	}
	roleID := u.getStringInt64(roleIDString)

	// 权限ID
	if _, ok := request.Post["permission_id"]; !ok {
		common.Error(err)
		return err
	}

	var permissionID []int64
	for _, value := range request.Post["permission_id"].Values {
		permissionID = append(permissionID, u.getStringInt64(value))
	}

	// 调用后端执行
	rsp, err := u.RoleService.UpdatePermission(ctx, &role.RolePermission{
		RoleId:       roleID,
		PermissionId: permissionID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(rsp)
	response.Body = string(bytes)
	return nil
}

// getPost 获取参数
func (u *UserApi) getPost(request *userApi.Request, key string) (string, error) {
	if _, ok := request.Post[key]; ok {
		return "", errors.New("参数异常")
	}
	return request.Post[key].Values[0], nil
}

// getStringInt64 字符串转int64
func (u *UserApi) getStringInt64(stringValue string) int64 {
	intValue, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		common.Error(err)
		return 0
	}
	return intValue
}

// getID 获取用户ID
func (u *UserApi) getID(request *userApi.Request) (int64, error) {
	// 检验
	if _, ok := request.Get["user_id"]; !ok {
		return 0, errors.New("参数异常")
	}

	// 获取ID进行转化
	idString := request.Get["user_id"].Values[0]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		common.Error(err)
		return 0, err
	}

	return id, nil
}

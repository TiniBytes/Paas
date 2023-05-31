package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"tini-paas/api/middlewareapi/proto/middlewareApi"
	"tini-paas/internal/middleware/proto/middleware"
	"tini-paas/pkg/common"
	"tini-paas/plugin/form"
)

// MiddlewareApi 中间件API处理（对API后端接口方法的实现）
type MiddlewareApi struct {
	MiddlewareService middleware.MiddlewareService
}

func (m *MiddlewareApi) AddMiddleware(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	addMiddleInfo := &middleware.MiddlewareInfo{}

	// 设置端口
	port, err := m.setMiddlePort(request)
	if err != nil {
		common.Error(err)
		return err
	}
	addMiddleInfo.MiddlePort = port

	// 设置环境变量
	addMiddleInfo.MiddleEnv = m.setMiddleEnv(request)

	// 设置存储
	addMiddleInfo.MiddleStorage = m.setMiddleStorage(request)

	// 获取类型
	middleTypeInfo := m.getMiddleType(request)

	// 判断不同的类型设置不同的值
	switch middleTypeInfo.MiddleTypeName {
	case "MYSQL":
		middleConfig := m.setMiddleConfig(request)
		addMiddleInfo.MiddleConfig = middleConfig
	}

	// 处理表单
	form.FormToMiddlewareStruct(request.Post, addMiddleInfo)

	// 调用后端执行添加操作
	rsp, err := m.MiddlewareService.AddMiddleware(ctx, addMiddleInfo)
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

func (m *MiddlewareApi) DeleteMiddleware(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	// 先获取middle_id
	if _, ok := request.Get["middle_id"]; !ok {
		response.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取middle_id
	middleIDString := request.Get["middle_id"].Values[0]
	ID, err := strconv.ParseInt(middleIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端执行删除
	rsp, err := m.MiddlewareService.DeleteMiddleware(ctx, &middleware.MiddlewareID{
		Id: ID,
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

func (m *MiddlewareApi) UpdateMiddleware(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	//TODO implement me
	panic("implement me")
}

func (m *MiddlewareApi) FindMiddlewareByID(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	// 先查询ID是否存在
	if _, ok := request.Get["middle_id"]; !ok {
		response.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取middle_id
	middleIDString := request.Get["middle_id"].Values[0]
	ID, err := strconv.ParseInt(middleIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 执行查询
	rsp, err := m.MiddlewareService.FindMiddlewareByID(ctx, &middleware.MiddlewareID{
		Id: ID,
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

func (m *MiddlewareApi) Call(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	allMiddleware, err := m.MiddlewareService.FindAllMiddleware(ctx, &middleware.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(allMiddleware)
	response.Body = string(bytes)
	return nil
}

func (m *MiddlewareApi) FindAllMiddlewareByTypeID(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	// 先查询type_id
	if _, ok := request.Get["middle_type_id"]; !ok {
		response.StatusCode = 200
		return errors.New("参数异常")
	}

	// 获取middle_type_id
	typeIDString := request.Get["middle_type_id"].Values[0]
	typeID, err := strconv.ParseInt(typeIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端服务执行查询
	rsp, err := m.MiddlewareService.FindAllMiddlewareByTypeID(ctx, &middleware.FindAllByTypeID{
		TypeId: typeID,
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

func (m *MiddlewareApi) AddMiddleType(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	typeInfo := &middleware.MiddleTypeInfo{
		MiddleTypeName:     request.Post["middle_type_name"].Values[0],
		MiddleTypeImageSrc: request.Post["middle_type_image_src"].Values[0],
		MiddleVersion: []*middleware.MiddleVersion{
			{
				MiddleDockerImage: "docker/consul",
				MiddleVersion:     "v1.0.1",
			},
		},
	}

	// 调用后端服务执行添加
	rsp, err := m.MiddlewareService.AddMiddleType(ctx, typeInfo)
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

func (m *MiddlewareApi) DeleteMiddleType(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	// 先查询middle_type_id
	if _, ok := request.Get["middle_type_id"]; !ok {
		response.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取middle_id
	middleTypeIDString := request.Get["middle_type_id"].Values[0]
	ID, err := strconv.ParseInt(middleTypeIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 调用后端执行删除
	rsp, err := m.MiddlewareService.DeleteMiddleType(ctx, &middleware.MiddleTypeID{
		Id: ID,
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

func (m *MiddlewareApi) UpdateMiddleType(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	//TODO implement me
	panic("implement me")
}

func (m *MiddlewareApi) FindMiddleTypeByID(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	// 先查询ID是否存在
	if _, ok := request.Get["middle_type_id"]; !ok {
		response.StatusCode = 500
		return errors.New("参数异常")
	}

	// 获取middle_id
	middleTypeIDString := request.Get["middle_type_id"].Values[0]
	ID, err := strconv.ParseInt(middleTypeIDString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 执行查询
	rsp, err := m.MiddlewareService.FindMiddleTypeByID(ctx, &middleware.MiddleTypeID{
		Id: ID,
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

func (m *MiddlewareApi) FindAllMiddleType(ctx context.Context, request *middlewareApi.Request, response *middlewareApi.Response) error {
	allMiddleType, err := m.MiddlewareService.FindAllMiddleType(ctx, &middleware.FindAll{})
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据回写
	response.StatusCode = 200
	bytes, _ := json.Marshal(allMiddleType)
	response.Body = string(bytes)
	return nil
}

// setMiddlePort 设置端口
func (m *MiddlewareApi) setMiddlePort(request *middlewareApi.Request) ([]*middleware.MiddlePort, error) {
	dataSlice, ok := request.Post["middle_port"]
	if !ok {
		return nil, errors.New("无端口")
	}

	// 特殊处理
	var middlePortSlice []*middleware.MiddlePort

	for _, value := range dataSlice.Values {
		parseInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			common.Error(err)
		}
		port := &middleware.MiddlePort{
			MiddlePort:     int32(parseInt),
			MiddleProtocol: "TCP",
		}

		middlePortSlice = append(middlePortSlice, port)
	}
	return middlePortSlice, nil
}

// setMiddleEnv 设置中间件环境变量
func (m *MiddlewareApi) setMiddleEnv(request *middlewareApi.Request) []*middleware.MiddleEnv {
	var envSlice []*middleware.MiddleEnv

	// 处理环境变量
	i := 1
	for {
		tag := "middle_env.key." + strconv.Itoa(i)
		valueTag := "middle_env.value." + strconv.Itoa(i)
		key, ok := request.Post[tag]
		if ok {
			env := &middleware.MiddleEnv{
				EbvKey:   key.Values[0],
				EnvValue: request.Post[valueTag].Values[0],
			}
			envSlice = append(envSlice, env)
		} else {
			break
		}
		i++
	}
	return envSlice
}

// setMiddleStorage 设置中间件存储
func (m *MiddlewareApi) setMiddleStorage(request *middlewareApi.Request) []*middleware.MiddleStorage {
	var storageSlice []*middleware.MiddleStorage

	// 处理存储
	i := 1
	for {
		nameTag := "middle_storage.name." + strconv.Itoa(i)
		sizeTag := "middle_storage.size." + strconv.Itoa(i)
		pathTag := "middle_storage.path." + strconv.Itoa(i)

		key, ok := request.Post[nameTag]
		if ok {
			value, _ := strconv.ParseFloat(request.Post[sizeTag].Values[0], 64)
			storage := &middleware.MiddleStorage{
				MiddleStorageName:       key.Values[0],
				MiddleStorageSize:       float32(value),
				MiddleStoragePath:       request.Post[pathTag].Values[0],
				MiddleStorageClass:      "csi-rbd-sc",
				MiddleStorageAccessMode: "ReadWriteOnce",
			}
			storageSlice = append(storageSlice, storage)
		} else {
			break
		}
		i++
	}
	return storageSlice
}

// getMiddleType 获取中间件类型
func (m *MiddlewareApi) getMiddleType(request *middlewareApi.Request) (middleTypeInfo middleware.MiddleTypeInfo) {
	typeValue, ok := request.Post["middle_type_id"]
	if !ok {
		return
	}

	typeID, err := strconv.ParseInt(typeValue.Values[0], 10, 64)
	if err != nil {
		// 转换失败
		common.Error(err)
		return
	}

	typeInfo, err := m.MiddlewareService.FindMiddleTypeByID(context.TODO(), &middleware.MiddleTypeID{
		Id: typeID,
	})
	if err != nil {
		common.Error(err)
		return
	}
	middleTypeInfo = *typeInfo
	return
}

// setMiddleConfig 设置中间件配置
func (m *MiddlewareApi) setMiddleConfig(request *middlewareApi.Request) *middleware.MiddleConfig {
	middleConfig := &middleware.MiddleConfig{}

	middleConfig.MiddleConfigRootUser = m.getValue(request, "middle_config_root_user")
	middleConfig.MiddleConfigRootPwd = m.getValue(request, "middle_config_root_pwd")
	middleConfig.MiddleConfigUser = m.getValue(request, "middle_config_user")
	middleConfig.MiddleConfigPwd = m.getValue(request, "middle_config_pwd")
	middleConfig.MiddleConfigDataBase = m.getValue(request, "middle_config_data_base")
	return middleConfig
}

// getValue 获取值
func (m *MiddlewareApi) getValue(request *middlewareApi.Request, key string) string {
	value, ok := request.Post[key]
	if ok {
		return value.Values[0]
	}
	return ""
}

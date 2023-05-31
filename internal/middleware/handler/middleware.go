package handler

import (
	"context"
	"strconv"
	"tini-paas/internal/middleware/model"
	"tini-paas/internal/middleware/proto/middleware"
	"tini-paas/internal/middleware/service"
	"tini-paas/pkg/common"
)

// MiddlewareHandler 中间件处理接口(对API后端接口方法的实现)
type MiddlewareHandler struct {
	// MiddlewareService 中间件操作接口
	MiddlewareService service.MiddlewareService

	// MiddleTypeService 中间件类型接口
	MiddleTypeService service.MiddleTypeService
}

func (m *MiddlewareHandler) AddMiddleware(ctx context.Context, info *middleware.MiddlewareInfo, response *middleware.Response) error {
	middleModel := &model.Middleware{}

	// 将info信息映射到middleModel
	err := common.SwapTo(info, middleModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 调用其它服务处理数据
	// 根据ID查询需要的镜像地址 -> 根据镜像地址创建
	imageAddress, err := m.MiddleTypeService.FindImageVersionByID(info.MiddleVersionId)
	if err != nil {
		common.Error(err)
		return err
	}
	info.MiddleDockerImageVersion = imageAddress

	// 在k8s中创建资源
	err = m.MiddlewareService.CreateToK8s(info)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 创建成功，写入数据库
	middleID, err := m.MiddlewareService.AddMiddleware(middleModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	response.Msg = "中间件创建成功，ID为：" + strconv.FormatInt(middleID, 10)
	common.Info(response.Msg)
	return nil
}

func (m *MiddlewareHandler) DeleteMiddleware(ctx context.Context, id *middleware.MiddlewareID, response *middleware.Response) error {
	// 先查询中间件信息
	middleModel, err := m.MiddlewareService.FindMiddlewareByID(id.Id)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 删除k8s中的资源
	err = m.MiddlewareService.DeleteFromK8s(middleModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	return nil
}

func (m *MiddlewareHandler) UpdateMiddleware(ctx context.Context, info *middleware.MiddlewareInfo, response *middleware.Response) error {
	err := m.MiddlewareService.UpdateToK8s(info)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 查询中间件信息
	middleModel, err := m.MiddlewareService.FindMiddlewareByID(info.Id)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 将info映射成middleModel
	err = common.SwapTo(info, middleModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 更新数据库信息
	err = m.MiddlewareService.UpdateMiddleware(middleModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	return nil
}

func (m *MiddlewareHandler) FindMiddlewareByID(ctx context.Context, id *middleware.MiddlewareID, info *middleware.MiddlewareInfo) error {
	middleModel, err := m.MiddlewareService.FindMiddlewareByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 转换数据格式
	err = common.SwapTo(middleModel, info)
	if err != nil {
		common.Error(err)
		return err
	}

	return nil
}

func (m *MiddlewareHandler) FindAllMiddleware(ctx context.Context, all *middleware.FindAll, rsp *middleware.AllMiddleware) error {
	allMiddleware, err := m.MiddlewareService.FindAllMiddleware()
	if err != nil {
		common.Error(err)
		return err
	}

	// 整理格式
	for _, middle := range allMiddleware {
		middleInfo := &middleware.MiddlewareInfo{}

		// 转换数据格式
		err = common.SwapTo(middle, middleInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		rsp.MiddlewareInfo = append(rsp.MiddlewareInfo, middleInfo)
	}
	return nil
}

func (m *MiddlewareHandler) FindAllMiddlewareByTypeID(ctx context.Context, id *middleware.FindAllByTypeID, rsp *middleware.AllMiddleware) error {
	allMiddleware, err := m.MiddlewareService.FindAllMiddlewareByTypeID(id.TypeId)
	if err != nil {
		common.Error(err)
		return err
	}

	// 整理格式
	for _, middle := range allMiddleware {
		middleInfo := &middleware.MiddlewareInfo{}

		// 转换数据格式
		err = common.SwapTo(middle, middleInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		rsp.MiddlewareInfo = append(rsp.MiddlewareInfo, middleInfo)
	}
	return nil
}

func (m *MiddlewareHandler) AddMiddleType(ctx context.Context, info *middleware.MiddleTypeInfo, response *middleware.Response) error {
	middleTypeModel := &model.MiddleType{}

	// 将info映射成middleTypeModel
	err := common.SwapTo(info, middleTypeModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 执行添加
	middleTypeID, err := m.MiddleTypeService.AddMiddleType(middleTypeModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	response.Msg = "中间件类型添加成功，ID为：" + strconv.FormatInt(middleTypeID, 10)
	common.Info(response.Msg)
	return nil
}

func (m *MiddlewareHandler) DeleteMiddleType(ctx context.Context, id *middleware.MiddleTypeID, response *middleware.Response) error {
	err := m.MiddleTypeService.DeleteMiddleType(id.Id)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	return nil
}

func (m *MiddlewareHandler) UpdateMiddleType(ctx context.Context, info *middleware.MiddleTypeInfo, response *middleware.Response) error {
	// 先查询中间件信息
	middleTypeModel, err := m.MiddleTypeService.FindMiddleTypeByID(info.Id)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 数据类型转换
	err = common.SwapTo(info, middleTypeModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}

	// 执行更新操作
	err = m.MiddleTypeService.UpdateMiddleType(middleTypeModel)
	if err != nil {
		common.Error(err)
		response.Msg = err.Error()
		return err
	}
	return nil
}

func (m *MiddlewareHandler) FindMiddleTypeByID(ctx context.Context, id *middleware.MiddleTypeID, info *middleware.MiddleTypeInfo) error {
	middleTypeModel, err := m.MiddleTypeService.FindMiddleTypeByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 数据类型转换
	err = common.SwapTo(middleTypeModel, info)
	if err != nil {
		common.Error(err)
		return err
	}

	return nil
}

func (m *MiddlewareHandler) FindAllMiddleType(ctx context.Context, all *middleware.FindAll, middleType *middleware.AllMiddleType) error {
	allMiddleType, err := m.MiddleTypeService.FindAllMiddleType()
	if err != nil {
		common.Error(err)
		return err
	}

	// 整理格式
	for _, middle := range allMiddleType {
		middleInfo := &middleware.MiddleTypeInfo{}

		// 数据类型转换
		err = common.SwapTo(middle, middleInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		middleType.MiddleTypeInfo = append(middleType.MiddleTypeInfo, middleInfo)
	}
	return nil
}

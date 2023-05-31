package handler

import (
	"context"
	"strconv"
	"tini-paas/internal/svc/model"
	"tini-paas/internal/svc/proto/svc"
	"tini-paas/internal/svc/service"
	"tini-paas/pkg/common"
)

// SvcHandler 服务处理
type SvcHandler struct {
	SvcService service.SvcService
}

// AddSvc 添加服务
func (s *SvcHandler) AddSvc(ctx context.Context, info *svc.SvcInfo, response *svc.Response) error {
	svcModel := &model.Svc{}

	// 将info数据类型转换为svcModel
	err := common.SwapTo(info, svcModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 在k8s中创建服务
	err = s.SvcService.CreateSvcToK8s(info)
	if err != nil {
		common.Error(err)
		return err
	}

	// 在k8s中创建成功 -> 在数据库中创建对应信息
	svcID, err := s.SvcService.AddSvc(svcModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 创建成功，回写数据
	common.Info("Svc 添加成功，ID为：" + strconv.FormatInt(svcID, 10))
	response.Msg = "Svc 添加成功，ID为：" + strconv.FormatInt(svcID, 10)
	return nil
}

// DeleteSvc 删除服务
func (s *SvcHandler) DeleteSvc(ctx context.Context, id *svc.SvcID, response *svc.Response) error {
	// 查询此ID的服务是否存在
	svcModel, err := s.SvcService.FindSvcByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 当前存在，可以在k8s中删除
	err = s.SvcService.DeleteFromK8s(svcModel)
	if err != nil {
		common.Error(err)
		return err
	}

	return nil
}

// UpdateSvc 更新服务
func (s *SvcHandler) UpdateSvc(ctx context.Context, info *svc.SvcInfo, response *svc.Response) error {
	// 先更新k8s数据
	err := s.SvcService.UpdateSvcToK8s(info)
	if err != nil {
		common.Error(err)
		return err
	}

	// 查询数据库中的svc
	svcModel, err := s.SvcService.FindSvcByID(info.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将新的info信息写入到svcModel
	err = common.SwapTo(info, svcModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将新组装好的svcModel更新到数据库
	return s.SvcService.UpdateSvc(svcModel)
}

// FindSvcByID 根据ID查找服务
func (s *SvcHandler) FindSvcByID(ctx context.Context, id *svc.SvcID, info *svc.SvcInfo) error {
	svcModel, err := s.SvcService.FindSvcByID(id.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将podModel信息包装成info
	err = common.SwapTo(svcModel, info)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

// FindAllSvc 查找全部服务
func (s *SvcHandler) FindAllSvc(ctx context.Context, all *svc.FindAll, rsp *svc.AllSvc) error {
	allSvc, err := s.SvcService.FindAllSvc()
	if err != nil {
		common.Error(err)
		return err
	}

	// 组装数据格式
	for _, m := range allSvc {
		svcInfo := &svc.SvcInfo{}

		err = common.SwapTo(m, svcInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		rsp.SvcInfo = append(rsp.SvcInfo, svcInfo)
	}
	return nil
}

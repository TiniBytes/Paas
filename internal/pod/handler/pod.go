package handler

import (
	"context"
	"strconv"
	"tini-paas/internal/pod/model"
	"tini-paas/internal/pod/proto/pod"
	"tini-paas/internal/pod/service"
	"tini-paas/pkg/common"
)

// PodHandler pod处理
type PodHandler struct {
	PodService service.PodService
}

// AddPod 添加pod
func (p *PodHandler) AddPod(ctx context.Context, info *pod.PodInfo, rsp *pod.Response) error {
	common.Info("添加pod")
	podModel := &model.Pod{}

	// 通过json tag 将info映射到podModel上
	err := common.SwapTo(info, podModel)
	if err != nil {
		common.Error(err)
		rsp.Msg = err.Error()
		return err
	}

	// 没有错误，执行创建逻辑
	err = p.PodService.CreateToK8s(info)
	if err != nil {
		common.Error(err)
		rsp.Msg = err.Error()
		return err
	}
	// 操作数据库
	podID, err := p.PodService.AddPod(podModel)
	if err != nil {
		common.Error(err)
		rsp.Msg = err.Error()
		return err
	}
	common.Info("Pod 添加成功，数据库ID为：" + strconv.FormatInt(podID, 10))
	rsp.Msg = "Pod 添加成功，数据库ID为：" + strconv.FormatInt(podID, 10)
	return nil
}

// DeletePod 删除pod
func (p *PodHandler) DeletePod(ctx context.Context, podID *pod.PodID, rsp *pod.Response) error {
	// 先在数据库中查找是否有此pod
	podModel, err := p.PodService.FindPodByID(podID.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 拿到podModel
	err = p.PodService.DeletedFromK8s(podModel)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

// FindPodByID 根据ID查找pod
func (p *PodHandler) FindPodByID(ctx context.Context, podID *pod.PodID, info *pod.PodInfo) error {
	// 先查询数据库pod信息
	podModel, err := p.PodService.FindPodByID(podID.Id)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将podModel信息包装成info
	err = common.SwapTo(podModel, info)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

// UpdatePod 更新pod
func (p *PodHandler) UpdatePod(ctx context.Context, info *pod.PodInfo, rsp *pod.Response) error {
	// 先更新k8s的pod
	err := p.PodService.UpdateToK8s(info)
	if err != nil {
		common.Error(err)
		return err
	}

	// 查询数据库中的pod
	podModel, err := p.PodService.FindPodByID(info.Id)
	if err != nil {
		common.Error(err)
		return err
	}
	// 将info信息更新到podModel
	err = common.SwapTo(info, podModel)
	if err != nil {
		common.Error(err)
		return err
	}

	// 将新组装好的podModel更新到数据库
	return p.PodService.UpdatePod(podModel)
}

// FindAllPod 查找全部pod
func (p *PodHandler) FindAllPod(ctx context.Context, findAll *pod.FindAll, allPod *pod.AllPod) error {
	// 先在数据库中查找全部pod信息
	pods, err := p.PodService.FindAllPod()
	if err != nil {
		common.Error(err)
		return err
	}

	// 包装成返回需要的格式
	for _, v := range pods {
		podInfo := &pod.PodInfo{}

		// 将每一个pod信息转为podInfo
		err = common.SwapTo(v, podInfo)
		if err != nil {
			common.Error(err)
			return err
		}

		allPod.PodInfo = append(allPod.PodInfo, podInfo)
	}
	return nil
}

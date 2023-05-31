package service

import (
	"context"
	"errors"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"tini-paas/internal/volume/model"
	"tini-paas/internal/volume/proto/volume"
	"tini-paas/internal/volume/repository"
	"tini-paas/pkg/common"
)

// VolumeService 存储卷接口
type VolumeService interface {
	AddVolume(*model.Volume) (int64, error)
	DeleteVolume(int64) error
	UpdateVolume(*model.Volume) error
	FindVolume(int64) (*model.Volume, error)
	FindAllVolume() ([]model.Volume, error)

	CreateVolumeToK8s(*volume.VolumeInfo) error
	DeleteVolumeFromK8s(*model.Volume) error
}

// NewVolumeService 初始化存储卷服务
func NewVolumeService(volumeRepository repository.VolumeRepository, clientSet *kubernetes.Clientset) VolumeService {
	return &VolumeDataService{
		VolumeRepository: volumeRepository,
		K8sClientSet:     clientSet,
		deployment:       &v1.Deployment{},
	}
}

// VolumeDataService 存储服务对象
type VolumeDataService struct {
	VolumeRepository repository.VolumeRepository
	K8sClientSet     *kubernetes.Clientset
	deployment       *v1.Deployment
}

func (v *VolumeDataService) AddVolume(volume *model.Volume) (int64, error) {
	return v.VolumeRepository.CreateVolume(volume)
}

func (v *VolumeDataService) DeleteVolume(i int64) error {
	return v.VolumeRepository.DeleteVolume(i)
}

func (v *VolumeDataService) UpdateVolume(volume *model.Volume) error {
	return v.VolumeRepository.UpdateVolume(volume)
}

func (v *VolumeDataService) FindVolume(i int64) (*model.Volume, error) {
	return v.VolumeRepository.FindVolumeByID(i)
}

func (v *VolumeDataService) FindAllVolume() ([]model.Volume, error) {
	return v.VolumeRepository.FindAll()
}

func (v *VolumeDataService) CreateVolumeToK8s(info *volume.VolumeInfo) error {
	volume := v.setVolume(info)

	// 先查询之前是否存在
	_, err := v.K8sClientSet.CoreV1().PersistentVolumeClaims(info.VolumeNamespace).Get(context.TODO(), info.VolumeName, v13.GetOptions{})
	if err != nil {
		// 之前不存在 -> 可以创建
		_, err = v.K8sClientSet.CoreV1().PersistentVolumeClaims(info.VolumeNamespace).Create(context.TODO(), volume, v13.CreateOptions{})
		if err != nil {
			// 创建失败
			common.Error(err)
			return err
		}
		common.Info("存储创建成功")
		return nil
	}

	// 之前存在
	common.Error("存储空间" + info.VolumeName + "已经存在")
	return errors.New("存储空间" + info.VolumeName + "已经存在")
}

func (v *VolumeDataService) DeleteVolumeFromK8s(volume *model.Volume) error {
	// 先从k8s中删除
	err := v.K8sClientSet.CoreV1().PersistentVolumeClaims(volume.VolumeNamespace).Delete(context.TODO(), volume.VolumeName, v13.DeleteOptions{})
	if err != nil {
		// 删除失败
		common.Error(err)
		return err
	}

	// 删除成功
	err = v.VolumeRepository.DeleteVolume(volume.ID)
	if err != nil {
		common.Error()
		return err
	}
	common.Info("删除存储ID" + strconv.FormatInt(volume.ID, 10) + "成功")
	return nil
}

// setVolume 设置pvc详情信息
func (v *VolumeDataService) setVolume(info *volume.VolumeInfo) *v12.PersistentVolumeClaim {
	pvc := &v12.PersistentVolumeClaim{}

	// 设置接口类型
	pvc.TypeMeta = v13.TypeMeta{
		Kind:       "PersistentVolumeClaim",
		APIVersion: "v1",
	}

	// 设置基础信息
	pvc.ObjectMeta = v13.ObjectMeta{
		Name:      info.VolumeName,
		Namespace: info.VolumeNamespace,
		Annotations: map[string]string{
			"pv.kubernetes.io/bound-by-controller":          "yes", // 绑定控制器自动绑定
			"volume.beta.kubernetes.io/storage-provisioner": "rbd.csi.ceph.com",
		},
	}

	// 设置存储动态信息
	pvc.Spec = v12.PersistentVolumeClaimSpec{
		AccessModes:      v.getAccessMode(info),
		Resources:        v.getResource(info),
		StorageClassName: &info.VolumeStorageClassName,
		VolumeMode:       v.getVolumeMode(info),
	}
	return pvc
}

// getAccessMode 获取访问模式
func (v *VolumeDataService) getAccessMode(info *volume.VolumeInfo) []v12.PersistentVolumeAccessMode {
	var pvm []v12.PersistentVolumeAccessMode

	var pm v12.PersistentVolumeAccessMode
	switch info.VolumeAccessMode {
	case "ReadWriteOnce":
		pm = v12.ReadWriteOnce
	case "ReadOnlyMany":
		pm = v12.ReadOnlyMany
	case "ReadWriteMany":
		pm = v12.ReadWriteMany
	case "ReadWriteOncePod":
		pm = v12.ReadWriteOncePod
	default:
		pm = v12.ReadWriteOnce
	}
	pvm = append(pvm, pm)
	return pvm
}

// getResource 获取资源配置
func (v *VolumeDataService) getResource(info *volume.VolumeInfo) v12.ResourceRequirements {
	source := v12.ResourceRequirements{}
	source.Requests = v12.ResourceList{
		"storage": resource.MustParse(strconv.FormatFloat(float64(info.VolumeRequest), 'f', 6, 64) + "Gi"),
	}
	return source
}

// getVolumeMode 获取存储类型
func (v *VolumeDataService) getVolumeMode(info *volume.VolumeInfo) *v12.PersistentVolumeMode {
	var pvm v12.PersistentVolumeMode

	switch info.VolumePersistentVolumeMode {
	case "Block":
		pvm = v12.PersistentVolumeBlock
	case "Filesystem":
		pvm = v12.PersistentVolumeFilesystem
	default:
		pvm = v12.PersistentVolumeFilesystem
	}
	return &pvm
}

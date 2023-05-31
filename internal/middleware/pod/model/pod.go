package model

/*
	Pod的状态
	1、挂起（Pending）：
		Pod 已被 Kubernetes 系统接受，但有一个或者多个容器镜像尚未创建。
		等待时间包括调度 Pod 的时间和通过网络下载镜像的时间，这可能需要花点时间。
	2、运行中（Running）：
		该 Pod 已经绑定到了一个节点上，Pod 中所有的容器都已被创建。
		少有一个容器正在运行，或者正处于启动或重启状态。
	3、成功（Succeeded）：
		Pod 中的所有容器都被成功终止，并且不会再重启。
	4、失败（Failed）：
		Pod 中的所有容器都已终止了，并且至少有一个容器是因为失败终止。
		也就是说，容器以非 0 状态退出或者被系统终止。
	5、未知（Unknown）：
		因为某些原因无法取得 Pod 的状态，通常是因为与 Pod 所在主机通信失败。
*/

// Pod pod属性
type Pod struct {
	// ID podApi id
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`

	// PodName pod名称
	PodName string `gorm:"unique_index;not_null" json:"pod_name"`

	// PodNamespace pod命名空间
	PodNamespace string `json:"pod_namespace"`

	// PodTeamID pod所属团队
	PodTeamID int64 `json:"pod_team_id"`

	// PodCpuMin pod使用cpu的最小值
	PodCpuMin float32 `json:"pod_cpu_min"`

	// PodCpuMax pod使用cpu的最大值
	PodCpuMax float32 `json:"pod_cpu_max"`

	// PodReplicas pod副本数量
	PodReplicas int32 `json:"pod_replicas"`

	// PodMemoryMin pod使用的内存最小值
	PodMemoryMin float32 `json:"pod_memory_min"`

	// PodMemoryMax pod使用的内存最大值
	PodMemoryMax float32 `json:"pod_memory_max"`

	// PodPort pod开放的端口
	PodPort []PodPort `gorm:"ForeignKey:PodID" json:"pod_port"`

	// PodEnv pod环境变量
	PodEnv []PodEnv `gorm:"ForeignKey:PodEnv" json:"pod_env"`

	// PodPullPolicy 镜像拉取策略
	// Always: 总是拉取
	// IfNotPresent: 默认值，本地有则使用本地镜像，不拉取
	// Never: 只使用本地镜像，从不拉取
	PodPullPolicy string `json:"pod_pull_policy"`

	// PodRestart pod重启策略
	// Always: 当容器失效时, 由kubelet自动重启该容器
	// OnFailure: 当容器终止运行且退出码不为0时, 由kubelet自动重启该容器
	// Never: 不论容器运行状态如何, kubelet都不会重启该容器
	PodRestart string `json:"pod_restart"`

	// PodType pod发布策略
	// 重建(recreate)：停止旧版本部署新版本
	// 滚动更新(rolling-update)：一个接一个地以滚动更新方式发布新版本
	// 蓝绿(blue/green)：新版本与旧版本一起存在，然后切换流量
	// 金丝雀(canary)：将新版本面向一部分用户发布，然后继续全量发布
	// A/B测(a/b testing)：以精确的方式（HTTP 头、cookie、权重等）向部分用户发布新版本。
	// A/B测实际上是一种基于数据统计做出业务决策的技术。
	// 在 Kubernetes 中并不原生支持，需要额外的一些高级组件来完成改设置（比如Istio、Linkerd、Traefik、或者自定义 Nginx/Haproxy 等）。
	// Recreate,Custom,Rolling
	PodType string `json:"pod_type"`

	// PodImage 使用的镜像名称
	PodImage string `json:"pod_image"`
}

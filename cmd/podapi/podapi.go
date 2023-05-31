package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/opentracing/opentracing-go"
	"net"
	"net/http"
	"strconv"
	"tini-paas/api/podapi/handler"
	podApi "tini-paas/api/podapi/proto/podApi"
	microPodService "tini-paas/internal/pod/proto/pod"
	"tini-paas/pkg/common"
	hystrix2 "tini-paas/plugin/hystrix"
)

var (
	hostIP               = "127.0.0.1"
	serviceHost          = hostIP // 服务地址
	servicePort          = "8082" // 服务端口
	consulHost           = hostIP // 注册中心地址
	consulPort     int64 = 8500   // 注册中心端口
	tracerHost           = hostIP // 链路追踪地址
	tracerPort           = 6831   // 链路追踪端口
	hystrixPort          = 9092   // 熔断端口（每个服务不能重复）
	prometheusPort       = 9192   // 监控端口（每个服务不能重复）
)

func main() {
	// 1、注册中心
	//newRegistry := consul.NewRegistry(func(options *registry.Options) {
	//	options.Addrs = []string{
	//		consulHost + ":" + strconv.FormatInt(consulPort, 10),
	//	}
	//})
	newRegistry := consul.NewRegistry(
		registry.Addrs(consulHost + ":" + strconv.FormatInt(consulPort, 10)),
	)
	fmt.Println(newRegistry.Options(), "2")

	// 2、链路追踪
	tracer, closer, err := common.NewTracer("go.micro.api.podApi", tracerHost+":"+strconv.Itoa(tracerPort))
	if err != nil {
		common.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 3、添加熔断器
	streamHandler := hystrix.NewStreamHandler()
	streamHandler.Start()

	// 4、添加日志，将日志采集到日志中心
	// 1) 需要程序日志打入到日志文件中
	// 2) 在程序中添加filebeat.yml 文件
	// 3) 启动filebeat, 启动命令 ./filebeat -e -c filebeat.yml
	fmt.Println("日志统一记录在根目录 micro.log 文件中，请点击查看日志")

	// 5、启动熔断监听程序
	go func() {
		//http://127.0.0.1:9092/turbine/turbine.stream
		//看板访问地址 http://127.0.0.1:9002/hystrix，url后面一定要带 /hystrix
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", strconv.Itoa(hystrixPort)), streamHandler)
		if err != nil {
			common.Error(err)
		}
	}()

	// 6、添加监控采集地址
	//http://127.0.0.1:9192/metrics
	common.PrometheusBoot(prometheusPort)

	// 7、创建API服务
	service := micro.NewService(
		// 自定义服务地址，必须写在其它参数前面
		micro.Server(server.NewServer(func(options *server.Options) {
			options.Advertise = serviceHost + ":" + servicePort
		})),
		micro.Name("go.micro.api.podApi"),
		micro.Version("latest"),
		// 指定服务端口
		micro.Address(":"+servicePort),
		// 添加注册中心
		micro.Registry(newRegistry),
		// 添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 作为客户端范围启动熔断
		micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// 添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
		// 添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	service.Init()

	// 调用pod微服务
	podService := microPodService.NewPodService("go.micro.service.pod", service.Client())
	// 注册控制器
	err = podApi.RegisterPodApiHandler(service.Server(), &handler.PodApi{PodService: podService})
	if err != nil {
		common.Error(err)
	}
	// 启动服务
	err = service.Run()
	if err != nil {
		common.Fatal(err)
	}
}

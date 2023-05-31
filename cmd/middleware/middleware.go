package main

import (
	"flag"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"tini-paas/config"
	"tini-paas/internal/middleware/handler"
	"tini-paas/internal/middleware/proto/middleware"
	"tini-paas/internal/middleware/repository"
	service2 "tini-paas/internal/middleware/service"
	"tini-paas/pkg/common"
	hystrix2 "tini-paas/plugin/hystrix"
)

var (
	hostIp               = "127.0.0.1" // 服务地址
	serviceHost          = hostIp      // 服务地址
	servicePort          = "8089"      // 服务端口
	consulHost           = hostIp      // 注册配置中心IP
	consulPort     int64 = 8500        // 注册配置中心端口
	tracerHost           = hostIp      // 链路追踪IP
	tracerPort           = 6831        // 链路追踪端口
	hystrixPort          = 9099        // 熔断器端口
	prometheusPort       = 9199        // 监控
)

func main() {
	// 1、注册中心
	newRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			consulHost + ":" + strconv.FormatInt(consulPort, 10),
		}
	})

	// 2、配置中心
	//consulConfig, err := common.GetConsulConfig(consulHost, consulPort, "/micro/consulConfig")
	//if err != nil {
	//	common.Error(err)
	//}

	// 3、使用配置中心连接MySQL
	//mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	mysqlInfo := config.GetMySQLConfig()
	// 初始化数据库
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@tcp("+mysqlInfo.Host+":"+mysqlInfo.Port+")/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	// 禁止复表
	db.SingularTable(true)

	// 4、添加链路追踪
	tracer, closer, err := common.NewTracer("go.micro.service.middleware", tracerHost+":"+strconv.Itoa(tracerPort))
	if err != nil {
		common.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 5、熔断器
	streamHandler := hystrix.NewStreamHandler()
	streamHandler.Start()
	// 添加监听程序
	go func() {
		//http://192.168.0.112:9092/turbine/turbine.stream
		//看板访问地址 http://127.0.0.1:9002/hystrix，url后面一定要带 /hystrix
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", strconv.Itoa(hystrixPort)), streamHandler)
		if err != nil {
			common.Error(err)
		}
	}()

	// 6、添加日志中心
	// 1) 需要程序日志打入到日志文件中
	// 2) 在程序中添加filebeat.yml 文件
	// 3) 启动filebeat, 启动命令 ./filebeat.yml -e -c filebeat.yml.yml
	fmt.Println("日志统一记录在根目录 micro.log 文件中，请点击查看日志")

	// 7、监控
	common.PrometheusBoot(prometheusPort)

	// 下载kubectl: https://kubernetes.io/docs/tasks/tools/#tabset-2
	// 1.curl.exe -LO "https://dl.k8s.io/v1.27.1/bin/windows/amd64/kubectl.exe.sha256"
	// 2.chmod +x ./kubectl
	// 3.sudo mv ./kubectl /usr/local/bin/kubectl
	// 4.sudo chown root: /usr/local/bin/kubectl
	// 5.kubectl version --client
	// 6.集群模式下直接拷贝服务端~/.kube/consulConfig 文件到本机 ~/.kube/confg 中
	//   注意：- config中的域名要能解析正确
	//        - 生产环境可以创建另一个证书
	// 7.kubectl get ns 查看是否正常
	//创建k8s连接
	//在集群外部使用
	// 将物理机config文件拷贝进docker
	//-v C:/Users/13158/.kube/consulConfig:/root/.kube/consulConfig
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "kubeConfig file 在当前系统的地址")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "kubeConfig file 在当前系统的地址")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		common.Fatal(err.Error())
	}

	//在集群中外的配置
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	panic(err.Error())
	//}

	// 创建程序可操作的客户端
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.Fatal(err.Error())
	}

	// 创建服务
	service := micro.NewService(
		// 自定义服务地址，且必须写在其它参数前面
		micro.Server(server.NewServer(func(options *server.Options) {
			options.Advertise = serviceHost + ":" + servicePort
		})),
		micro.Name("go.micro.service.middleware"),
		micro.Version("latest"),
		// 指定服务端口
		micro.Address(":"+servicePort),
		// 添加注册中心
		micro.Registry(newRegistry),
		// 添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 添加熔断，作为客户端使用
		micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// 添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
	)

	// 初始化服务
	service.Init()

	// 初始化数据表
	//err = repository.NewSvcRepository(db).InitTable()
	//if err != nil {
	//	common.Fatal(err)
	//}

	// 注册句柄
	middlewareService := service2.NewMiddlewareService(repository.NewMiddlewareRepository(db), clientSet)
	middleTypeService := service2.NewMiddleTypeService(repository.NewMiddleTypeRepository(db))
	err = middleware.RegisterMiddlewareHandler(service.Server(), &handler.MiddlewareHandler{
		// 注册两个服务接口
		MiddlewareService: middlewareService,
		MiddleTypeService: middleTypeService,
	})
	if err != nil {
		return
	}

	// 启动服务
	err = service.Run()
	if err != nil {
		common.Fatal(err)
	}
}

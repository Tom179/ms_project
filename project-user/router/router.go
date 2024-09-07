package router

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	user_grpc "test.com/project-grpc/user/login"
	"test.com/project-user/config"
	loginServiceV1 "test.com/project-user/pkg/service/login.service.v1"
)

/*
type Router interface { //路由的统一规范和接口，需实现
	SetRoute(r *gin.Engine)
}
*/
/*
type RegisterRouter struct { //某个类
}

	func NewRegistRouter() *RegisterRouter {
		return &RegisterRouter{}
	}

func (*RegisterRouter) Route(ro Router, r *gin.Engine) { //调用接口中的路由实现方法

		ro.Route(r)
	}
*/
//var routers []Router

type gRPCconfig struct { //这个类用来表示一个grpc微服务模块
	Addr         string
	RegisterFunc func(server *grpc.Server)
}

func RegistGrpc() *grpc.Server {
	ggg := gRPCconfig{Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			user_grpc.RegisterLoginServiceServer(g, loginServiceV1.New()) //将自定义的微服务结构体注册到grpcServer中
		}} //定义方法,未调用

	s := grpc.NewServer() //创建grpc服务端，也就是上面说的grpcServer
	ggg.RegisterFunc(s)
	lis, err := net.Listen("tcp", ggg.Addr)
	if err != nil {
		log.Println("cannot listen")
	}

	go func() {
		err = s.Serve(lis) //启动微服务
		if err != nil {
			log.Println("启动服务器失败", err)
			return
		}
	}()
	return s
}

func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG) //传入etcd地址
	resolver.Register(etcdRegister)
	//这个resolver拿来干嘛？

	info := discovery.Server{ //传入服务地址
		Name:    config.C.GC.Name,    //user
		Addr:    config.C.GC.Addr,    //127.0.0.1:8881
		Version: config.C.GC.Version, //1.0.0
		Weight:  config.C.GC.Weight,
	}

	r := discovery.NewRegister(config.C.EtcdConfig.Addrs, logs.LG) //创建注册器
	_, err := r.RegistService(info, 2)                             //注册服务
	if err != nil {
		log.Fatalln(err)
	}
}

package router

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	project "test.com/project-grpc/project"
	project_service_v1 "test.com/project-project/pkg/service/project.service.v1"
	"test.com/project-user/config"
)

type gRPCconfig struct { //这个类用来表示一个grpc微服务模块
	Addr         string
	RegisterFunc func(server *grpc.Server)
}

func RegistGrpc() *grpc.Server {
	ggg := gRPCconfig{Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			project.RegisterProjectServiceServer(g, project_service_v1.New())
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

	info := discovery.Server{ //project模块识别了user的配置..
		Name:    config.C.GC.Name,
		Addr:    config.C.GC.Addr,
		Version: config.C.GC.Version,
		Weight:  config.C.GC.Weight,
	}

	r := discovery.NewRegister(config.C.EtcdConfig.Addrs, logs.LG) //创建注册器
	_, err := r.RegistService(info, 2)                             //注册服务
	if err != nil {
		log.Fatalln(err)
	}
}

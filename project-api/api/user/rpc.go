package user

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"test.com/project-api/config"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	"test.com/project-grpc/user/login"
)

var LoginServiceClient login.LoginServiceClient //引用这个客户端可以远程调用微服务

// 调用grpc服务
func InitRpcUserClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister) //注册解析器构造器

	///////版本1.0.0
	//这个函数是grpc连接，并不实际的grpc调用
	conn, err := grpc.Dial("etcd://user", grpc.WithTransportCredentials(insecure.NewCredentials())) //地址直接填的名字，由resolver来解析，resolver里builder实现的etcd方式
	//自动调用builder接口中的Scheme函数和build函数，build函数通过target.URL.path和target.URL.host【也就是服务名和服务地址】来获取etcd的服务信息。

	if err != nil {
		log.Fatalf("连接server失败:%v", err)
	}

	LoginServiceClient = login.NewLoginServiceClient(conn) //创建grpc客户端
}

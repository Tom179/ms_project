package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"test.com/project-api/config"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	project "test.com/project-grpc/project"
)

var ProjectServiceClient project.ProjectServiceClient //引用这个客户端可以远程调用微服务

func InitRpcProjectClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd://project", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接server失败:%v", err)
	}

	ProjectServiceClient = project.NewProjectServiceClient(conn)
}

package main

import (
	"github.com/gin-gonic/gin"
	common "test.com/project-common"
	"test.com/project-project/config"
	"test.com/project-project/router"
)

func main() {
	r := gin.Default()
	gc := router.RegistGrpc()   //
	router.RegisterEtcdServer() //grpc服务注册到etcd

	stop := func() {
		gc.Stop()
	} //停止函数

	common.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}

package main

import (
	"github.com/gin-gonic/gin"
	common "test.com/project-common"
	"test.com/project-user/config"
	_ "test.com/project-user/internal/dao"
	_ "test.com/project-user/internal/database/gorm"
	"test.com/project-user/router"
)

func main() {
	r := gin.Default()

	//router.InitRouter(r)
	gc := router.RegistGrpc()   //注册grpc
	router.RegisterEtcdServer() //grpc服务注册到etcd
	stop := func() {
		gc.Stop()
	} //停止函数

	common.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)

	//r.Run(":8080")
	//fmt.Println("后续操作")
}

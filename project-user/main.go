package main

import (
	"github.com/gin-gonic/gin"
	common "test.com/project-common"
	_ "test.com/project-user/api" //什么意思，导入
	"test.com/project-user/config"
	"test.com/project-user/router"
)

func main() {
	r := gin.Default()

	router.InitRouter(r)
	gc := router.RegistGrpc() //注册grpc
	stop := func() {
		gc.Stop()
	} //停止函数

	common.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)

	//r.Run(":8080")
	//fmt.Println("后续操作")
}

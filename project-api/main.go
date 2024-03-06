package main

import (
	"github.com/gin-gonic/gin"
	"test.com/project-api/config"
	"test.com/project-api/router"

	_ "test.com/project-api/api" //执行了api-user-router中的init函数,所以一来就添加了user服务
	common "test.com/project-common"
)

func main() {
	r := gin.Default()
	//路由
	router.InitRouter(r)

	common.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)

	//r.Run(":80")
	//fmt.Println("后续操作")
}

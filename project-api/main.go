package main

import (
	"github.com/gin-gonic/gin"
	_ "test.com/project-api/api" //执行了api-user-router中的init函数,所以一来就添加了user服务
	"test.com/project-api/config"
	"test.com/project-api/router"
	common "test.com/project-common"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	common.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)
}

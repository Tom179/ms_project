package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/router"
)

func init() {
	log.Println("初始化API模块中的project路由")
	ru := &RouterProject{}
	router.AddRoute(ru)
}

type RouterProject struct {
}

func (*RouterProject) SetRoute(r *gin.Engine) { //实现路由接口
	InitRpcProjectClient() //连接grpc服务
	h := NewProjectHandler()
	r.POST("/project/index", h.index) //获取验证码

}

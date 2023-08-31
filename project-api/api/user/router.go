package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/router"
)

func init() {
	log.Println("初始化API模块中的user路由")
	ru := &RouterUser{}
	router.AddRoute(ru)
}

type RouterUser struct {
}

func (*RouterUser) SetRoute(r *gin.Engine) { //实现路由接口
	InitRpcUserClient() //连接grpc服务
	h := NewUserHandler()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}

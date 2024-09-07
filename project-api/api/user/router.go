package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/api/MiddleWare"
	"test.com/project-api/api/rpc"
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
	rpc.InitRpcUserClient() //连接grpc服务
	h := NewUserHandler()
	r.POST("/project/login/getCaptcha", h.getCaptcha) //获取验证码
	r.POST("/project/login/register", h.register)     //注册
	r.POST("/project/login", h.login)
	//r.POST("/project/index", h.index)

	MessageGroup := r.Group("project/organization")
	MessageGroup.Use(MiddleWare.TokenVerify)
	MessageGroup.POST("/_getOrgList", h.MyOrg)
}

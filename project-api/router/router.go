package router

import (
	"github.com/gin-gonic/gin"
)

type Router interface { //路由的统一规范和接口，需实现
	SetRoute(r *gin.Engine)
}

/*
type RegisterRouter struct { //某个类
}

	func NewRegistRouter() *RegisterRouter {
		return &RegisterRouter{}
	}

func (*RegisterRouter) Route(ro Router, r *gin.Engine) { //调用接口中的路由实现方法

		ro.Route(r)
	}
*/
var routers []Router

func InitRouter(r *gin.Engine) {
	/*regestRuoter := NewRegistRouter()
	regestRuoter.Route(&user.UserRouter{}, r)*/
	for _, ro := range routers {
		ro.SetRoute(r)
	}

}

func AddRoute(ro ...Router) { //添加路由
	routers = append(routers, ro...)
}

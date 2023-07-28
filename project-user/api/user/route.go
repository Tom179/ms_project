package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"test.com/project-user/router"
)

func init() {
	fmt.Println("初始化添加路由")
	router.AddRoute(&UserRouter{}) //添加路由实现
}

// 路由的实现
type UserRouter struct {
}

func (*UserRouter) SetRoute(r *gin.Engine) {
	hander := NewUserHandler()                             //注意是通过指针来引用的
	r.POST("/project/login/getCaptcha", hander.GetCaptcha) //
}

package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	"test.com/project-user/config"
	loginServiceV1 "test.com/project-user/pkg/service/login.service.v1"
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

type gRPCconfig struct {
	Addr         string
	RegisterFunc func(server *grpc.Server)
}

func RegistGrpc() *grpc.Server {
	c := gRPCconfig{Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			loginServiceV1.RegisterLoginServiceServer(g, loginServiceV1.NewLoginService()) //将自定义的Server注册到grpcServer中
		}}
	s := grpc.NewServer()
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		log.Println("cannot listen")
	}

	go func() {
		err = s.Serve(lis) //启动服务器
		if err != nil {
			log.Println("启动服务器失败", err)
			return
		}
	}()

	return s
}

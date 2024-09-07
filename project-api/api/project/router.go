package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/api/MiddleWare"
	"test.com/project-api/api/rpc"
	"test.com/project-api/router"
)

func init() {
	log.Println("初始化API模块中的project路由")
	rp := &RouterProject{}
	router.AddRoute(rp)
}

type RouterProject struct {
}

func (*RouterProject) SetRoute(r *gin.Engine) {
	rpc.InitRpcProjectClient()
	h := NewProjectHandler()

	projectGroup := r.Group("/project/index")
	projectGroup.Use(MiddleWare.TokenVerify)
	projectGroup.POST("", h.index)

	projectGroup1 := r.Group("/project/project")
	projectGroup1.Use(MiddleWare.TokenVerify)
	projectGroup1.POST("/selfList", h.MyProject)
	projectGroup1.POST("/", h.MyProject)

}

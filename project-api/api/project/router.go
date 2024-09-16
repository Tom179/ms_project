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

	projectGroup := r.Group("/project")
	projectGroup.Use(MiddleWare.TokenVerify)
	projectGroup.POST("/index", h.index)

	projectGroup.POST("/project/selfList", h.MyProject)
	projectGroup.POST("/project", h.MyProject)
	projectGroup.POST("/project_template", h.ProjectTemplate)
	projectGroup.POST("/project/save", h.SaveProject)
	projectGroup.POST("/project/read", h.ReadProject)

}

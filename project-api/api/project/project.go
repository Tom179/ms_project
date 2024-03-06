package project

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/project"
)

type ProjectHandler struct {
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
}

func (*ProjectHandler) index(c *gin.Context) {
	result := &common.Result{}
	req := IndexReq{}
	c.ShouldBindJSON(&req)
	fmt.Println(req.Token)

	rsp, err := ProjectServiceClient.Index(context.Background(), &project.IndexMessage{Token: req.Token})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success(rsp.Menus))
}

type IndexReq struct {
	Token string `json:"token"`
}

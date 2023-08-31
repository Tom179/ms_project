package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	common "test.com/project-common"
	"test.com/project-common/errs"
	service "test.com/project-user/pkg/service/login.service.v1"
	"time"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (*UserHandler) getCaptcha(c *gin.Context) { //路由映射到此方法
	result := &common.Result{}
	mobile := c.PostForm("mobile")
	//fmt.Println("mobile", mobile)
	//👇发起grpc调用（前提是已经将loginServiceClient实例化）
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rsp, err := LoginServiceClient.GetCaptcha(ctx, &service.CaptchaRequest{Mobile: mobile})
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //从错误中解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Code))
}

package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"test.com/project-api/pkg/model/user"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/user/login"
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
	rsp, err := LoginServiceClient.GetCaptcha(ctx, &login.CaptchaRequest{Mobile: mobile})
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //从错误中解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, result.Success(rsp.Code))
}

func (*UserHandler) register(c *gin.Context) {

	result := &common.Result{}
	var req user.RegisterReq
	err := c.ShouldBind(&req) //获取请求参数
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
	}

	if err := req.Verify(); err != nil { //验证格式
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	msg := &login.RegisterRequest{}
	err = copier.Copy(msg, req) //用工具库给msg赋值
	//fmt.Println("copy的msg为", msg) ////
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "结构体复制错误"))
	}

	_, err = LoginServiceClient.Register(ctx, msg)

	//gRPC调用
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return //出现异常，直接返回
	}

	c.JSON(http.StatusOK, result.Success(""))
}

package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"test.com/project-api/api/rpc"
	"test.com/project-api/pkg/model/user"
	common "test.com/project-common"
	"test.com/project-common/errs"
	login "test.com/project-grpc/user/login"
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
	//👇开启grpc链接，前提是已经将loginServiceClient实例化）
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rsp, err := rpc.LoginServiceClient.GetCaptcha(ctx, &login.CaptchaRequest{Mobile: mobile}) //grpc返回的err，和go中自带的err类型不一样，要手动解析grpcErr中的信息
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

	fmt.Println("准备grpc调用")
	fmt.Println("前端传给后端的password为:", msg.Password)
	_, err = rpc.LoginServiceClient.Register(ctx, msg) //这才是具体的grpc调用啊
	fmt.Println("接收到的grpc调用的返回值err为: ", err)

	//gRPC调用
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		fmt.Println("grpc客户端接收到的code和msg分别为：", code, " ", msg)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return //出现异常，直接返回
	}

	c.JSON(http.StatusOK, result.Success(""))
}

func (*UserHandler) login(c *gin.Context) {
	result := &common.Result{}
	req := LoginReq{}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	fmt.Println("输入的用户名:", req.Account, " 密码:", req.Password)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	msg := &login.LoginMessage{}
	err = copier.Copy(msg, req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy有误"))
		return
	}

	//grpc调用
	loginResp_grpc, err := rpc.LoginServiceClient.Login(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	fmt.Println("【grpc调用的响应为】:")
	fmt.Println(loginResp_grpc)

	rsp := LoginRsp{}

	err = copier.Copy(&rsp, loginResp_grpc)

	fmt.Println("【api返回的Web响应为】:")
	fmt.Println(rsp)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy出错"))
	}

	c.JSON(http.StatusOK, result.Success(rsp))
}

type MyOrgReq struct {
	Id int `json:"id"`
}

func (*UserHandler) MyOrg(c *gin.Context) {
	result := common.Result{}
	req := MyOrgReq{}

	c.ShouldBindJSON(&req)
	lc := rpc.LoginServiceClient
	myorgMsg := login.MyOrgReqGrpc{}
	copier.Copy(&myorgMsg, req)

	rsp, err := lc.MyOrganization(context.Background(), &myorgMsg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success(rsp))

}

type LoginReq struct {
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
}

type LoginRsp struct {
	Member           Member             `json:"member"`
	TokenList        TokenList          `json:"tokenList"`
	OrganizationList []OrganizationList `json:"organizationList"`
}
type Member struct {
	Id               int64  `json:"id"`
	Code             string `json:"code"` //对id进行加密，可解密
	Name             string `json:"name"`
	Mobile           string `json:"mobile"`
	Status           int    `json:"status"`
	CreateTime       string `json:"create_time"`
	LastLoginTime    string `json:"last_login_time"`
	OrganizationCode string `json:"organization_code"`
}

type TokenList struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	TokenType      string `json:"tokenType"`
	AccessTokenExp int64  `json:"accessTokenExp"`
}

type OrganizationList struct {
	//Id          int64  `json:"id"`
	Code        string `json:"code"` //对id进行加密
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	OwnerCode   string `json:"owner_code"`
	CreateTime  string `json:"create_time"`
	Personal    int32  `json:"personal"`
	Address     string `json:"address"`
	Province    int32  `json:"province"`
	City        int32  `json:"city"`
	Area        int32  `json:"area"`
}

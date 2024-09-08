package MiddleWare

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"test.com/project-api/api/rpc"
	common "test.com/project-common"
	"test.com/project-common/errs"
	login "test.com/project-grpc/user/login"
	"time"
)

// 1.从header中获取token
// 2.调用User服务进行token认证
// 3. 处理结果，将信息放入gin的上下文
func TokenVerify(c *gin.Context) {
	result := &common.Result{}

	token := c.GetHeader("Authorization")
	if token == "" {
		fmt.Println("token为空")
		c.Abort()
	}

	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "")
	}
	fmt.Println("Authorization头中得到,替换后的token为:", token)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	resp, err := rpc.LoginServiceClient.TokenVerify(ctx, &login.LoginMessage{Token: token})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		c.Abort()
		return
	}
	c.Set("memberId", resp.Member.Id)
	c.Set("memberName", resp.Member.Name)
	c.Set("organizationCode", resp.Member.OrganizationCode)
	c.Next()
}

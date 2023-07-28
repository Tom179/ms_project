package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	common "test.com/project-common"
	"test.com/project-user/pkg/dao"
	"test.com/project-user/pkg/model"
	"test.com/project-user/pkg/repo"
	"time"
)

type UserHandler struct {
	cache repo.Cache //存储接口
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		cache: dao.Rc,
	}
}

func (this *UserHandler) GetCaptcha(c *gin.Context) {
	rsp := &common.Result{}
	phoneNumber := c.PostForm("mobile")
	if !common.VerifyMobile(phoneNumber) {
		c.JSON(http.StatusOK, rsp.Fail(model.IllegalMobile, "手机号不合法"))
		return

	}
	code := RandomCaptCha()
	go func() { //发送短信
		time.Sleep(2 * time.Second)
		zap.L().Info("api发送短信 info") //假设发送成功
		//logs.LG.Debug("api发送短信 debug")
		//zap.L().Error("api发送短信 error")

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := this.cache.Put(ctx, "REGISTER:"+phoneNumber, code, 15*time.Minute) //存入redis
		if err != nil {
			log.Printf("验证码存入redis出错,cause by: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, rsp.Success(code))
}

func RandomCaptCha() string {
	// 使用时间作为种子，确保每次调用结果不一样
	rand.Seed(time.Now().UnixNano())
	// 生成6位随机数，范围在[100000, 999999]
	return strconv.Itoa(rand.Intn(900000) + 100000)
}

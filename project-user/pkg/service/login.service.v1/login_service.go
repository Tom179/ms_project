package login_service_v1

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"strconv"
	common "test.com/project-common"
	"test.com/project-user/pkg/dao"
	"test.com/project-user/pkg/repo"
	"time"
)

type LoginService struct {
	UnimplementedLoginServiceServer            //实现grpc
	cache                           repo.Cache //自己新增参数
}

func NewLoginService() *LoginService {
	return &LoginService{cache: dao.Rc}
}

func (lg *LoginService) GetCaptcha(c context.Context, msg *CaptchaRequest) (*CaptchaResponse, error) {
	mobile := msg.Mobile

	if !common.VerifyMobile(mobile) {
		return nil, errors.New("手机号不合法")
	}
	code := RandomCaptCha()
	go func() { //发送短信
		time.Sleep(2 * time.Second)
		zap.L().Info("api发送短信 info") //假设发送成功
		//logs.LG.Debug("api发送短信 debug")
		//zap.L().Error("api发送短信 error")

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := lg.cache.Put(ctx, "REGISTER:"+mobile, code, 15*time.Minute) //存入redis
		if err != nil {
			log.Printf("验证码存入redis出错,cause by: %v\n", err)
		}
	}()

	return &CaptchaResponse{}, nil
}

func RandomCaptCha() string {
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机数种子

	min := 100000
	max := 999999
	randomNumber := rand.Intn(max-min+1) + min

	return strconv.Itoa(randomNumber)
}

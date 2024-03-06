package model

import (
	"test.com/project-common/errs"
)

/*const (//直接定义未int型
	IllegalMobile common.ResponseCode = 2001 //手机号不合法
)*/

var ( //直接定义错误
	RedisError        = errs.NewError(999, "redis错误")
	DBError           = errs.NewError(998, "db错误")
	IllegalMobile     = errs.NewError(10102001, "手机号不合法")
	CaptchaNotExist   = errs.NewError(10102002, "验证码不存在")
	InCorrectCaptcha  = errs.NewError(10102002, "验证码不正确")
	EmailExisted      = errs.NewError(10102003, "email已经存在")
	AccountExisted    = errs.NewError(10102004, "账号已经存在")
	MobileExisted     = errs.NewError(10102005, "手机号已经存在")
	AccuntAndPwdError = errs.NewError(10102007, "账号密码不正确")
)

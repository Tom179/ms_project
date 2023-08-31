package model

import (
	"test.com/project-common/errs"
)

/*const (//直接定义未int型
	IllegalMobile common.ResponseCode = 2001 //手机号不合法
)*/

var ( //直接定义错误
	IllegalMobile = errs.NewError(2001, "手机号不合法")
)

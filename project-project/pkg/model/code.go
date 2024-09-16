package model

import (
	"test.com/project-common/errs"
)

var ( //直接定义错误
	RedisError = errs.NewError(999, "redis错误")
	DBError    = errs.NewError(998, "db错误")
)

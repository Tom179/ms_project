package errs

import "fmt"

type ErrorCode int
type BError struct { //自定义错误
	Code ErrorCode
	Msg  string
}

func (e *BError) Error() string {
	return fmt.Sprintf("code:%v,msg:%s", e.Code, e.Msg)
}

func NewError(code ErrorCode, msg string) *BError { //构造函数
	return &BError{
		Code: code,
		Msg:  msg,
	}
}

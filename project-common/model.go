package common //响应格式

type ResponseCode int
type Result struct {
	Code ResponseCode `json:"code"` //json标签字段，让结构体无论是从前端接收还是返回给前端。前端的键名都是json指定的格式：通常为小写
	Msg  string       `json:"msg"`
	Data any          `json:"data"`
}

func (r *Result) Success(data any) *Result { //赋值并返回success
	r.Code = 200
	r.Msg = "成功！"
	r.Data = data
	return r
}

func (r *Result) Fail(code ResponseCode, msg string) *Result {
	r.Code = code
	r.Msg = msg
	return r
}

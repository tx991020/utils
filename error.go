package utils

//自定义错误和返回格式
type oerror struct {
	Code    int `json:"code"`
	Message string `json:"message"`
}

type OniError *oerror

func NewError(code int, Message string) *oerror {
	return &oerror{Code: code, Message: Message}
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

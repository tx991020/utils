package utils

//自定义错误和返回格式
type oerror struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type OniError *oerror

func NewError(code int, Message string) *oerror {
	return &oerror{Code: code, Message: Message}
}

type Response struct {
	Code int         `json:"errcode"`
	Msg  string      `json:"errmsg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

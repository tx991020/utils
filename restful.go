package utils




type hterror struct {
	Code int    `json:"errcode,omitempty"`
	Msg  string `json:"errmsg,omitempty"`
}



func NewError(code int, msg string) *hterror{
	return &hterror{Code: code, Msg: msg}
}


type Response struct {
	Code int         `json:"errcode"`
	Msg  string      `json:"errmsg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

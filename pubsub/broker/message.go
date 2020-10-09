package broker



// 客户端读写消息
type message struct {
	// websocket.TextMessage 消息类型
	messageType int
	data        []byte
}

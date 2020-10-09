


package pubsub

// 订阅查询过滤器
type QueryFilter interface {
	GetTopic() string
}

// 订阅到的消息
type SubMessage interface {
	GetConent() []byte
}

type Subscription interface {
	// Unsubscribe cancels the sending of events to the data channel
	// and closes the error channel.
	Unsubscribe()
	// Err returns the subscription error channel. The error channel receives
	// a value if there is an issue with the subscription (e.g. the network connection
	// delivering the events has been closed). Only one value will ever be sent.
	// The error channel is closed by Unsubscribe.
	Err() <-chan error
}

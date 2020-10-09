package client

import (
	"context"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/tx991020/utils/pubsub"
)

type client struct {
	err  <-chan error
	conn *websocket.Conn
}

func Dial(addr string) (*client, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return nil, err
	}

	cli := &client{}
	cli.err = make(<-chan error)
	cli.conn = c

	return cli, nil
}

func (c *client) Subscribe(cxt context.Context, filter pubsub.QueryFilter,
	ch chan<- pubsub.SubMessage) (pubsub.Subscription, error) {

	// 向服务端发送订阅主题的消息
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(filter.GetTopic())); err != nil {
		return nil, err
	}

	// 从管道中读取数据
	go func() {
		if err := c.readLoop(ch); err != nil {
			return
		}
	}()

	return c, nil
}

func (c *client) readLoop(ch chan<- pubsub.SubMessage) error {

	for {
		msgType, msgData, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		msg := &message{
			msgType,
			msgData,
		}

		ch <- msg
	}

	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Unsubscribe() {
	return
}

func (c *client) Err() <-chan error {
	return c.err
}

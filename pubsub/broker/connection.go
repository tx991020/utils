// Copyright (c) 2019. sample-go Inc. All rights reserved.

package broker

import (
	"errors"
	"regexp"
	"sync"
	"time"

	"github.com/tx991020/utils/logger"

	"github.com/gorilla/websocket"
)

var (
	errConnClose = errors.New("conntion 连接已关闭")
)

type connection struct {
	broker     *Broker
	channel    *channel
	clientConn *websocket.Conn

	inChan  chan *message
	outChan chan *message

	mutex     sync.Mutex
	isClosed  bool
	closeChan chan byte
	id        int64
}

func (conn *connection) processLoop() {
	conn.clientConn.SetReadLimit(maxMessageSize)
	conn.clientConn.SetReadDeadline(time.Now().Add(pongWait))

	for {
		// 读一个message
		_, data, err := conn.clientConn.ReadMessage()
		if err != nil {

			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			logger.Errorf("消息读取出现错误:%v", err.Error())
			if conn.channel != nil {
				conn.channel.removeConn(conn)
			}
			conn.close()
			conn.broker.updateConnTotal(-1)
			return
		}

		channelName := string(data)
		conn.broker.subscribe(conn, channelName)

		select {
		case <-conn.closeChan:
			return
		}
	}
}

// 发送消息给客户端
func (conn *connection) writeLoop() {

	ticker := time.NewTicker(pingPeriod)

	defer ticker.Stop()

	for {
		select {
		case msg := <-conn.outChan:
			if err := conn.clientConn.WriteMessage(msg.messageType, msg.data); err != nil {
				logger.Errorf("发送消息给客户端发生错误:%v", err.Error())
				conn.close()
				return
			}
		case <-conn.closeChan:
			return
		case <-ticker.C:
			conn.clientConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.clientConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 写入消息到队列中
func (conn *connection) write(messageType int, data []byte) error {
	select {
	case conn.outChan <- &message{messageType, data}:
	case <-conn.closeChan:
		return errConnClose
	}
	return nil
}

// 检查channelName的合法性
func (conn *connection) checkChannelName(name string) bool {
	match, _ := regexp.MatchString("^[0-9a-zA-Z=]{2,64}$", name)
	return match
}

// 关闭连接
func (conn *connection) close() {
	conn.clientConn.Close()
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if conn.isClosed == false {
		conn.isClosed = true
		close(conn.closeChan)
	}
}

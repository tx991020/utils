// Copyright (c) 2019. pubsub Inc. All rights reserved.

package broker

import (

	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type waitGroupWrapper struct {
	sync.WaitGroup
}

func (w *waitGroupWrapper) wrap(cb func(argvs ...interface{}), argvs ...interface{}) {
	w.Add(1)
	go func() {
		cb(argvs...)
		w.Done()
	}()
}

type channel struct {
	sync.RWMutex
	name         string
	connections  map[int64]*connection
	waitGroup    *waitGroupWrapper
	messageCount uint64
	exitFlag     int32

}

func NewChannel(name string) *channel {
	c := &channel{}
	c.name = name
	c.connections = make(map[int64]*connection)
	c.waitGroup = &waitGroupWrapper{}

	return c
}

// 添加连接
func (c *channel) addConn(conn *connection) {
	c.Lock()
	c.connections[conn.id] = conn
	c.Unlock()
}

// 删除连接
func (c *channel) removeConn(conn *connection) {
	c.Lock()
	delete(c.connections, conn.id)
	c.Unlock()
}

// 渠道下发消息
func (c *channel) notify(message []byte) {
	c.RLock()
	defer c.RUnlock()

	for _, conn := range c.connections {
		c.waitGroup.wrap(func(args ...interface{}) {
			conn := args[0].(*connection)
			message := args[1].([]byte)
			conn.write(websocket.TextMessage, message)
		}, conn, message)
	}
}

// 等待消息发送完成
func (c *channel) wait() {
	c.waitGroup.Wait()
}

func (c *channel) exiting() bool {
	return atomic.LoadInt32(&c.exitFlag) == 1
}

func (c *channel) exit() {
	if !atomic.CompareAndSwapInt32(&c.exitFlag, 0, 1) {
		return
	}
	c.wait()
}

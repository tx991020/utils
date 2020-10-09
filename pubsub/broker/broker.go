package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"github.com/gorilla/websocket"
)

const (

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time allowd to write the message to the peer.
	writeWait = (pingPeriod * 11) / 10

	// Maximum message size (unit:byte) allowed from peer.
	maxMessageSize = 1024
)

var (
	ErrUpgrade               = errors.New("升级为websocket失败")
	ErrNotFoundChannelFormat = "没有找到channel:%s"
)

type Broker struct {
	sync.RWMutex
	ctx        context.Context
	server     *http.Server
	maxConn    int64
	maxConnId  *int64
	connTotal  *int64
	log        log.Logger
	addr       string
	wsupgrader *websocket.Upgrader
	channels   map[string]*channel
}

func NewBroker(ctx context.Context, addr string, maxConn int64) *Broker {
	b := &Broker{}
	b.maxConn = maxConn
	b.maxConnId = new(int64)
	b.connTotal = new(int64)
	b.wsupgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	b.channels = make(map[string]*channel)
	b.addr = addr
	b.ctx = ctx

	return b
}

func (b *Broker) Run() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", b.handler)

	b.server = &http.Server{
		Addr:    b.addr,
		Handler: mux,
	}

	return b.server.ListenAndServe()
}

// 对连接进行计数
func (b *Broker) updateConnTotal(delta int64) int64 {
	return atomic.AddInt64(b.connTotal, delta)
}

// 连接的流水号
func (b *Broker) incrMaxconnId() int64 {
	return atomic.AddInt64(b.maxConnId, 1)
}

// 判断连接是否已满
func (b *Broker) connNotFull() bool {
	return *b.connTotal < b.maxConn
}

// 处理连接
func (b *Broker) handler(resp http.ResponseWriter, req *http.Request) {
	var err error
	clientConn, err := b.wsupgrader.Upgrade(resp, req, nil)
	if err != nil {
		b.log.Fatalf("升级为websocket失败:%v", err.Error())
		return
	}

	if !b.connNotFull() {
		// 超过最大连接数
		clientConn.WriteMessage(websocket.TextMessage, []byte("连接超过最大限制"))
		clientConn.Close()
		return
	}

	b.incrMaxconnId()

	conn := &connection{
		broker:     b,
		clientConn: clientConn,
		inChan:     make(chan *message, 1000),
		outChan:    make(chan *message, 1000),
		closeChan:  make(chan byte),
		isClosed:   false,
		id:         *b.maxConnId,
	}

	b.updateConnTotal(1)

	go conn.processLoop()
	go conn.writeLoop()

}

// 订阅
func (b *Broker) subscribe(conn *connection, channelName string) {

	b.RLock()
	ch, found := b.channels[channelName]
	b.RUnlock()

	if !found {
		ch = NewChannel(channelName)
		conn.channel = ch
		ch.addConn(conn)
		b.Lock()
		b.channels[channelName] = ch
		b.Unlock()
	} else {
		conn.channel = ch
		ch.addConn(conn)
	}
}

// 取消订阅
func (b *Broker) unsubscribe(conn *connection, channelName string) {

	b.RLock()
	ch, found := b.channels[channelName]
	b.RUnlock()

	if found {
		ch.removeConn(conn)
		ch.exit()
	}
}

// 添加通道
func (b *Broker) AddChannel(channelName string) (bool, error) {

	b.RLock()
	ch, found := b.channels[channelName]
	b.RUnlock()

	if !found {
		ch = NewChannel( channelName)
		b.Lock()
		b.channels[channelName] = ch
		b.Unlock()
	}

	return true, nil
}

// 广播消息
func (b *Broker) Publish(channelName string, message []byte) (bool, error) {

	b.RLock()
	ch, found := b.channels[channelName]
	b.RUnlock()

	if !found {
		return false, fmt.Errorf(ErrNotFoundChannelFormat, channelName)
	}

	ch.notify(message)
	ch.wait()

	return true, nil
}

// 关闭server
func (b *Broker) Shutdown() error {
	return b.server.Shutdown(b.ctx)
}

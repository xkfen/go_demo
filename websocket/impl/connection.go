package impl

import (
	"github.com/gorilla/websocket"
	"sync"
	"github.com/pkg/errors"
)

type WebConnection struct {
	Conn *websocket.Conn
	// 输入的message
	InChan chan []byte
	// 输出的message
	OutChan chan []byte
	// chan读写的时候都会阻塞
	CloseChan chan byte

	Mutex sync.Mutex
	IsClosed bool
}

// 封装websocket长连接
func InitConnection(wsConn *websocket.Conn) (conn *WebConnection, err error) {
	conn = &WebConnection{
		Conn:      wsConn,
		InChan:    make(chan []byte, 1000),
		OutChan:   make(chan []byte, 1000),
		CloseChan: make(chan byte, 1),
	}

	// 在初始化连接的时候启动goroutine
	go conn.ReadLoop()
	go conn.WriteLoop()

	return conn, nil
}

// 读消息
func (conn *WebConnection) ReadMessage() (data []byte, err error) {
	select{
	case data= <- conn.InChan:
	case <-conn.CloseChan:
		err = errors.New("connection is closed")
	}
	data = <-conn.InChan
	return
}

// 写消息
func (conn *WebConnection) WriteMessage(data []byte) (err error) {
	select{
	case conn.OutChan <- data:
	case conn.CloseChan:
		err = errors.New("connection is closed")
	}

	return nil
}

func (conn *WebConnection) Close() {
	// websocket的close是线程安全的，可重入的close
	conn.Conn.Close()

	// 只执行一次,保证channel智慧执行一次
	conn.Mutex.Lock()
	if !conn.IsClosed {
		close(conn.CloseChan)
		conn.IsClosed = true
	}
	conn.Mutex.Unlock()
}

func (conn *WebConnection) ReadLoop() {
	for {
		_, data, err := conn.Conn.ReadMessage()
		if err != nil {
			goto ERR
		}
		// 这里会发生阻塞，等待inChan有空闲的位置
		select {
		// 如果有数据
		case conn.InChan <- data:
			// 如果已经被关闭
		case <-conn.CloseChan:
			goto ERR

		}

	}
ERR:
	conn.Close()
}

// 发送消息
func (conn *WebConnection) WriteLoop() {
	var data []byte
	for {

		select {
		case data = <-conn.OutChan:
		case <-conn.CloseChan:
			goto ERR
		}

		if err := conn.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}

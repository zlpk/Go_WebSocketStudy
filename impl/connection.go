package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte
	isClose bool
	mutex sync.Mutex
}

// 连接初始化封装 API，外部可调用
func InitConnection(wsConn *websocket.Conn)(conn *Connection, err error)  {
	conn = &Connection{
		wsConn:wsConn,
		inChan:make(chan []byte, 1000),
		outChan:make(chan []byte, 1000),
		closeChan:make(chan byte, 1),
	}

	// 启动读协程，如果有消息就放到 inChan 队列中
	go conn.readLoop()

	// 启动写协程，从 outChan 中读取数据，如果有数据就发出去
	go conn.writeLoop()

	return
}

// 读消息封装 API，外部可调用
func (conn *Connection) ReadMessage()(data []byte, err error)  {
	select {
	case data = <- conn.inChan:
	case <- conn.closeChan:
		err = errors.New("connection closed")
	}
	return
}

// 写消息封装 API，外部可调用
func (conn *Connection) WriteMessage(data []byte)(err error)  {
	select {
	case conn.outChan <- data:
	case <- conn.closeChan:
		err = errors.New("connection closed")
	}
	return
}

// 关闭连接封装 API，外部可调用
func (conn *Connection) Close() {

	// 线程安全，可重入的 Close
	conn.wsConn.Close()

	// 为保证 Close 的线程安全，关闭通道之前先加锁
	conn.mutex.Lock()

	// 通道只能关闭一次
	if !conn.isClose {
		close(conn.closeChan)
		conn.isClose = true
	}

	conn.mutex.Unlock()
}

// 内部实现，循环接收消息，放入 inChan
func (conn *Connection) readLoop() {
	var (
		data []byte
		err error
	)

	for {
		if _,data,err = conn.wsConn.ReadMessage();err != nil {
			goto ERR
		}

		// 为了防止网络异常关闭连接后，一直阻塞在 inChan，使用 Select 语句进行规避
		select {
		case conn.inChan <- data:
		case <- conn.closeChan:
			goto ERR
		}
	}

ERR:
	conn.Close()
}

// 内部实现，循环取 inChan 数据，发送出去
func (conn *Connection) writeLoop() {
	var (
		data []byte
		err error
	)

	for {
		// 为了防止网络异常关闭连接后，一直阻塞在 outChan，使用 Select 语句进行规避
		select {
		case data = <- conn.outChan:
		case <- conn.closeChan:
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data);err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
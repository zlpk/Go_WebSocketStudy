package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"./impl"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

/*
 * 处理 http 请求 handler
 */
func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		data []byte
		err    error
		conn *impl.Connection
	)

	// http 升级为 websocket
	if wsConn,err = upgrader.Upgrade(w, r, nil); err != nil {
		fmt.Println("建立 WebSocket 连接失败!")
		return
	} else {
		fmt.Println("建立 WebSocket 连接成功!")
	}

	// 封装 websocket
	if conn,err = impl.InitConnection(wsConn);err != nil {
		goto ERR
	}

	// 心跳协程，每 1 秒给客户端发送一个心跳包，同时也说明封装是线程安全的
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heart beat"));err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// 循环将客户端发来的消息发送回去
	for {
		if data,err = conn.ReadMessage();err != nil {
			goto ERR
		}

		fmt.Printf("接收客户端消息：%s\n", data)

		if err = conn.WriteMessage(data);err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
	fmt.Println("连接关闭！")
}

func main()  {
	// 建路由
	http.HandleFunc("/ws", wsHandler)

	// 指定地址监听并提供服务
	http.ListenAndServe("0.0.0.0:8080", nil)
}

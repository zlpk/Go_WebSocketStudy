package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
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
		err    error
		data   []byte
	)

	// http 升级为 websocket
	if wsConn,err = upgrader.Upgrade(w, r, nil); err != nil {
		fmt.Println("建立 WebSocket 连接失败!")
		return
	} else {
		fmt.Println("建立 WebSocket 连接成功!")
	}

	// 循环接收消息并将接收到的消息原样发回去
	for {
		if _,data,err = wsConn.ReadMessage(); err != nil {
			goto ERR
		} else {
			fmt.Printf("接收客户端消息：%s\n", data)
		}

		if err = wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

// 错误处理
ERR:
	fmt.Println("关闭 WebSocket 连接!")
	wsConn.Close()
}

func main()  {
	// 建路由
	http.HandleFunc("/ws", wsHandler)

	// 指定地址监听并提供服务
	http.ListenAndServe("0.0.0.0:8080", nil)
}

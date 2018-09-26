package main

import "net/http"

/*
处理 http 请求 handler
 */
func wsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello ZLPK!</h1>"))
}

func main()  {
	// 建路由
	http.HandleFunc("/ws", wsHandler)

	// 指定地址监听并提供服务
	http.ListenAndServe("0.0.0.0:8080", nil)
}

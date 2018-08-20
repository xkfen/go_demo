package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"github.com/go-kit/kit/util/conn"
	"go_demo/websocket/impl"
	"time"
)

var (
	Upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request){
	// Upgrade:websocket  握手
	wsConn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn, err := impl.InitConnection(wsConn)
	if err != nil {
		fmt.Println(err.Error())
		goto ERR
	}

	// 启动协程，每隔1秒给客户端发送心跳
	go func() {
		for {
			if err := conn.WriteMessage([]byte("heart beat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()


	for{
		data, err := conn.ReadMessage()
		if err != nil {
			goto ERR
		}
		if err := conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}
	ERR:
		conn.Close()
}


func main(){
	// http标准库配置路由
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
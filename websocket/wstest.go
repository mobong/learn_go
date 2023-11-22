package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type ReqMsg struct {
	Id int `json:"id"`
}

func WsClient() {
	var dialer *websocket.Dialer
	//通过Dialer连接websocket服务器
	conn, _, err := dialer.Dial("ws://127.0.0.1:8001", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	go do(conn)
	send(conn, "test", &ReqMsg{Id: 1})
	time.Sleep(10 * time.Second)
}

func do(conn *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%s\n", err)
		}
	}()
	for {
		_, message, _ := conn.ReadMessage()
		msg := &RspBody{}
		data, _ := UnZip(message)
		if err := json.Unmarshal(data, msg); err == nil {
			fmt.Printf("received: %s, code:%d, %v\n", msg.Name, msg.Code, msg.Msg)
		} else {
			fmt.Println("json Unmarshal error:", err)
		}
	}
}

func send(conn *websocket.Conn, name string, Msg interface{}) {
	msg := &ReqBody{Name: name, Msg: Msg}

	if data, err := json.Marshal(msg); err == nil {
		data, _ = Zip(data)
		conn.WriteMessage(websocket.BinaryMessage, data)
	}
}

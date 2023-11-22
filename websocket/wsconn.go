package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type ServerConn struct {
	wsSocket *websocket.Conn
	outChan  chan *WsMsgRsp // 写队列
}

type ReqBody struct {
	Seq  int64       `json:"seq"`
	Name string      `json:"name"`
	Msg  interface{} `json:"msg"`
}

type RspBody struct {
	Seq  int64       `json:"seq"`
	Name string      `json:"name"`
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

type WsMsgReq struct {
	Body *ReqBody
}

type WsMsgRsp struct {
	Body *RspBody
}

type RspMsg struct {
	Id    int    `json:"id"`
	World string `json:"world"`
}

func NewServerConn(wsSocket *websocket.Conn) *ServerConn {
	return &ServerConn{
		wsSocket: wsSocket,
		outChan:  make(chan *WsMsgRsp, 1000),
	}
}

func (sc *ServerConn) Start() {
	go sc.wsReadLoop()
	go sc.wsWriteLoop()
}

func (sc *ServerConn) wsReadLoop() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("wsReadLoop error", err)
			sc.Close()
		}
	}()

	for {
		_, data, err := sc.wsSocket.ReadMessage()
		if err != nil {
			break
		}
		data, err = UnZip(data)
		if err != nil {
			fmt.Println("wsReadLoop Unzip error:", err)
			continue
		}
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			fmt.Println("json Unmarshal error:", err)
		} else {
			fmt.Println(*body)
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: body.Seq, Msg: &RspMsg{
				Id:    1,
				World: "hello",
			}}}
			sc.outChan <- rsp
		}
	}
	sc.Close()
}

func (sc *ServerConn) wsWriteLoop() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("wsWriteLoop error", err)
			sc.Close()
		}
	}()

	for {
		select {
		// 取一个消息
		case msg := <-sc.outChan:
			// 写给websocket
			sc.write(msg.Body)
		}
	}
}

func (sc *ServerConn) Close() {
	sc.wsSocket.Close()
}

func (sc *ServerConn) write(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data, err = Zip(data)
	if err != nil {
		return err
	}
	err = sc.wsSocket.WriteMessage(websocket.BinaryMessage, data)
	return err
}

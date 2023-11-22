package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"learn_go/conf"
	"net/http"
	"strconv"
)

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type server struct {
	addr string
}

func NewServer() *server {
	conf.ConfInit()
	host := viper.GetString("webserver.host")
	port := viper.GetInt("webserver.port")
	addr := host + ":" + strconv.Itoa(port)
	return &server{addr: addr}
}

func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	http.ListenAndServe(s.addr, nil)
}

func (s *server) wsHandler(resp http.ResponseWriter, req *http.Request) {
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}
	conn := NewServerConn(wsSocket)
	conn.Start()
}

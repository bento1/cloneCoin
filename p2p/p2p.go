package p2p

import (
	"fmt"
	"net/http"

	"github.com/github.com/bento1/cloneCoin/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Port 3000 will upgrade the request from 4000
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	openPort := utils.Splitter(r.URL.Query().Get("openPort"), ":", 1)
	request := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && request != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil) // 3000이 4000으로 보내는 conn
	utils.HandleErr(err)
	initPeer(conn, request, openPort)
}

func AddPeer(address, port, openPort string) { //이함수는 3000이 요청하게 된다.
	// :4000 is requesting upgrade :3000
	fmt.Printf("ws://%s:%s/ws\n", address, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil) //upgrade가 완료되면 4000이 3000으로 보내는 conn
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	sendNewestBlock(peer)

}

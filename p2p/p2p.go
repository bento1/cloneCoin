package p2p

import (
	"fmt"
	"net/http"

	"github.com/github.com/bento1/cloneCoin/blockchain"
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
	fmt.Printf("%s want to upgrade to ws\n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil) // 3000이 4000으로 보내는 conn
	utils.HandleErr(err)
	initPeer(conn, request, openPort)
}

func AddPeer(address, port, openPort string, broadcast bool) { //이함수는 3000이 요청하게 된다.
	// :4000 is requesting upgrade :3000
	fmt.Printf("%s want to connect to %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil) //upgrade가 완료되면 4000이 3000으로 보내는 conn
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	if broadcast {
		BroadcastNewPeer(peer)
		return
	}
	sendNewestBlock(peer)
	// 자신이 가진 커넥션의 주소로 추가를 알려야함

}

func BroadcastNewBlock(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, peer := range Peers.value {
		notifyNewBlock(b, peer)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.value {
		notifyNewTx(tx, p)
	}
}

func BroadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.value {
		//새로 추가된 peer에서는 메세지를 보내면 안됨
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}

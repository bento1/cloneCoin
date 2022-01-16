package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peers struct {
	value map[string]*peer
	m     sync.Mutex
}

var Peers = peers{
	value: make(map[string]*peer),
}

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	key     string
	address string
	port    string
}

func ALLPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()
	// /peers
	// {
	// 	add:port :{}
	// }
	// [add:port,add:port,add:port,add:port,add:port]
	var keys []string
	for key := range p.value {
		keys = append(keys, key)
	}
	return keys
}
func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.value, p.key) // 데이터 레이스 만듬
}
func (p *peer) read() {
	// delete peer in case of err
	defer p.close() // a무한 루프니깐 // 데이터 레이스 만듬

	for {
		m := Message{}
		err := p.conn.ReadJSON(&m) //들어올떄까지 잡고 있다.block operate 자동으로 byte 를 Json으로 마샬파싱
		if err != nil {
			break
		}
		handleMsg(&m, p)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		message, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, message)
	}

}
func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)

	p := peer{
		conn:    conn,
		inbox:   make(chan []byte),
		address: address,
		key:     key,
		port:    port,
	}

	go p.read()           // 3000이 웹소켓으로 업그레이드하면서
	go p.write()          //4000 요청일 경우 3000이 가진 최근 블럭 정보를 받아옴
	Peers.value[key] = &p //여기서 data race를 야기할 수 있다. 껏다가 킬때

	return &p
}

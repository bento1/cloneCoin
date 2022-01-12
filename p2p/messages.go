package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/github.com/bento1/cloneCoin/blockchain"
	"github.com/github.com/bento1/cloneCoin/utils"
)

type MessageKind int

const (
	// MessageNewestBlock       MessageKind = 1
	// MessageAllBlocksResqust  MessageKind = 2
	// MessageAllBlocksResponse MessageKind = 3 iota쓰면 자동으로 value랑 type을 지정해준다
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksResqust
	MessageAllBlocksResponse
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJson(payload),
	}
	mJson := utils.ToJson(m)
	return mJson
}
func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, b) //kind를 줌
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Println(payload)
	}
	// fmt.Printf("Peer : %s, Sent a message with kind of : %d\n", p.key, m.Kind)
}

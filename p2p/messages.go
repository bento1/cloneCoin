package p2p

import (
	"encoding/json"

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
func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResqust, nil)
	p.inbox <- m

}
func sendAllBlock(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.BlockChain()))
	p.inbox <- m
}
func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock: //3000이 요청함
		var payload blockchain.Block // 4000의 체인
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		//fmt.Println(payload)
		b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			// request all the bloicks form
			requestAllBlocks(p)
		} else {
			//send  our blocks to 4000
			sendNewestBlock(p)
		}
	case MessageAllBlocksResqust:
		sendAllBlock(p)
	case MessageAllBlocksResponse:
		var payload []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))

	}
	// fmt.Printf("Peer : %s, Sent a message with kind of : %d\n", p.key, m.Kind)
}

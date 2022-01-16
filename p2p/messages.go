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
	MessageNewBlockNotify
	MessageNewTxNotify
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
	fmt.Printf("Send newestblock to  %s\n", p.key)
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
func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}
func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}
func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock: //3000이 요청함
		fmt.Printf("Receive newestblock from  %s\n", p.key) // read 한다. sendNewestBlock에서 MessageNewestBlock 를 보냈기 떄문에 이부분을 Read에서 Json parse하여 알게됨
		var payload blockchain.Block                        // 4000의 체인
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		//fmt.Println(payload)
		b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			// request all the bloicks form
			fmt.Printf("Requset all blocks to  %s\n", p.key)
			requestAllBlocks(p)
		} else {
			//send  our blocks to 4000
			fmt.Printf("Send newestblock form  %s\n", p.key)
			sendNewestBlock(p)
		}
	case MessageAllBlocksResqust:
		fmt.Printf("%s wants all blocks\n", p.key)
		sendAllBlock(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Receive all blocks from %s\n", p.key)
		var payload []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.BlockChain().Replace(payload)
	case MessageNewBlockNotify:
		fmt.Printf("broadcast new blocks from %s\n", p.key)
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.BlockChain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		fmt.Printf("broadcast new transaction from %s\n", p.key)
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	}

}

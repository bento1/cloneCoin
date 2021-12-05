package blockchain

import (
	"sync"

	"github.com/bento1/cloneCoin/db"
	"github.com/bento1/cloneCoin/utils"
)

type blockchain struct { //이제 마지막 해쉬만, 길이가 몇인지만 알면된다.
	NewestHash string `json:"newesthash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func BlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}
func (b *blockchain) persist() {
	db.SaveBlockChain(utils.ToBytea(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()

}

package blockchain

import (
	"crypto/sha256"
	"fmt"

	"github.com/github.com/bento1/cloneCoin/db"
	"github.com/github.com/bento1/cloneCoin/utils"
)

type Block struct {
	Data         string `json:"data"`
	Hash         string `json:"hash"`
	PreviousHash string `json:"previoushash,omitempty"`
	Height       int    `json:"height"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytea(b))
}
func createBlock(data string, previoushash string, height int) *Block {
	block := Block{
		Data:         data,
		Hash:         "",
		PreviousHash: previoushash,
		Height:       height,
	}
	payload := block.Data + block.PreviousHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return &block
}

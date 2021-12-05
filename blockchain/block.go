package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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

func (b *Block) toBytea() []byte {
	// return []byte(*b)// 이렇게 안된다.
	//gob이라는 패키지를 사용
	var blockbuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockbuffer)
	utils.HandleErr(encoder.Encode(b))
	return blockbuffer.Bytes()
}
func (b *Block) persist() {
	db.SaveBlock(b.Hash, b.toBytea())
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

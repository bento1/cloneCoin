package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/github.com/bento1/cloneCoin/db"
	"github.com/github.com/bento1/cloneCoin/utils"
)

// const difficulty int = 2

type Block struct {
	Hash         string `json:"hash"`
	PreviousHash string `json:"previoushash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"` // 유저가 채굴할떄 쓰는 변경가능한 값
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

func persistBlock(b *Block) {
	db.SaveBlock(b.Hash, utils.ToBytea(b))
}
func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		fmt.Println("target : ", target)
		fmt.Println("nonce : ", b.Nonce)
		fmt.Println("hash : ", hash)

		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		}
		b.Nonce++
	}

}
func createBlock(previoushash string, height int, difficulty int) *Block {
	block := Block{
		Hash:         "",
		PreviousHash: previoushash,
		Height:       height,
		Difficulty:   difficulty,
		Nonce:        0,
	} //통쨰로
	block.mine()
	block.Transactions = Mempool().txToConfirm() //채굴이 끝난 후 진행되어야함
	persistBlock(&block)
	return &block
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytea(b, data)
}
func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

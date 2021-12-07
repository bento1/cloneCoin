package blockchain

import (
	"fmt"
	"sync"

	"github.com/github.com/bento1/cloneCoin/db"
	"github.com/github.com/bento1/cloneCoin/utils"
)

type blockchain struct { //이제 마지막 해쉬만, 길이가 몇인지만 알면된다.
	NewestHash        string `json:"newesthash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentdifficulty"`
}

var b *blockchain
var once sync.Once

const defaultDifficulty int = 2
const difficultyInterval int = 5
const blockInterval int = 2

func (b *blockchain) restore(data []byte) {
	utils.FromBytea(b, data)
}

func BlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0, defaultDifficulty}
			fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
			// search checkpoint onthe db
			// restore b from bytea
			checkpoint := db.CheckPoint()
			if checkpoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}

		})
	}
	fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
	return b
}
func (b *blockchain) persist() {
	db.SaveBlockChain(utils.ToBytea(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()

}

func (b *blockchain) Blocks() []*Block {
	//previous hash를 계속 호출한다.
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PreviousHash != "" {
			hashCursor = block.PreviousHash
		} else {
			break
		}
	}
	return blocks
}
func (b *blockchain) recalculateDifficulty() int {
	//최근 difficultyinterval동안 timestamp를 알아본다.
	//difficultyinterval * blockinterval 안쪽이면 쉽게
	//
	allblock := b.Blocks()
	newestBlock := allblock[0]
	recalculateBlock := allblock[difficultyInterval-1]
	takenTime := (newestBlock.Timestamp / 60) - (recalculateBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if takenTime < expectedTime {
		return b.CurrentDifficulty + 1

	} else if takenTime > expectedTime {
		return b.CurrentDifficulty - 1
	} else {
		return b.CurrentDifficulty
	}
}
func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		//recalculate difficulty
		return b.recalculateDifficulty()
	} else {
		//아니면 이전 difficulty를 불러옴
		return b.CurrentDifficulty
	}
}

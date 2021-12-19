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

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}
func FindTx(b *blockchain, targetId string) *Tx {
	for _, tx := range Txs(b) {
		if tx.Id == targetId {
			return tx
		}
	}
	return nil
}
func restoreBlockchain(b *blockchain, data []byte) {
	utils.FromBytea(b, data)
}

func BlockChain() *blockchain {
	once.Do(func() {
		b = &blockchain{"", 0, defaultDifficulty}
		fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
		// search checkpoint onthe db
		// restore b from bytea
		checkpoint := db.CheckPoint()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			fmt.Println("Restoring...")
			restoreBlockchain(b, checkpoint)
		}

	})
	fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
	return b
}
func persistBlockchain(b *blockchain) {
	db.SaveBlockChain(utils.ToBytea(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, difficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)

}

func Blocks(b *blockchain) []*Block {
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
func recalculateDifficulty(b *blockchain) int {
	//최근 difficultyinterval동안 timestamp를 알아본다.
	//difficultyinterval * blockinterval 안쪽이면 쉽게
	//
	allblock := Blocks(b)
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
func difficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		//recalculate difficulty
		return recalculateDifficulty(b)
	} else {
		//아니면 이전 difficulty를 불러옴
		return b.CurrentDifficulty
	}
}

func UTxOutsByAddress(b *blockchain, address string) []*UTxOut {
	var UTxOuts []*UTxOut
	createSTxOut := make(map[string]bool)
	// 여기에 있으면output에 걸러야함
	// 거를 리스트를 저장한다.
	//map으로
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {

					createSTxOut[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := createSTxOut[tx.Id]; !ok {
						//존재하지 않는다 사용한적이없다.
						//현재코드는 이미 mempool에 올라간 utxout을 검사하지 않는다.
						//block에 추가 되기 전까지 할당된 utxout은 또다시 사용되면 안된다.

						// UTxOuts = append(UTxOuts, &UTxOut{tx.Id, index, output.Amount})

						UTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(UTxOut) {
							UTxOuts = append(UTxOuts, UTxOut)
						}

					}
				}

			}
		}
	}
	return UTxOuts
}
func BalanceByAddress(b *blockchain, address string) int {
	txOuts := UTxOutsByAddress(b, address)
	var amount int
	for _, txout := range txOuts {
		amount += txout.Amount
	}
	return amount
}

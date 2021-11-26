package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

//ver1 21-11-25 only blockchain develop
//여러 기능을 추가하면서 refactoring을 계속 할 것임 ... transaction, database..
type block struct {
	data         string //transaction 등이 바뀔수 있다.
	hash         string
	previousHash string
}
type blockchain struct {
	// blocks []block
	blocks []*block // 복사하고싶지 않음
}

var b *blockchain //singleton pattern 이녀석을  외부에서 읽게함-> 1개의 인스턴스만 존재하게됨
var once sync.Once

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}
func (b *block) calculateHash() {
	hash := sha256.Sum256([]byte(b.data + b.previousHash))
	b.hash = fmt.Sprintf("%x", hash)
}
func GetBlockChain() *blockchain { //singleton pattern 인스턴스를 외부 읽는 메소드-> 1개의 인스턴스만 존재하게됨
	if b == nil {
		// b = &blockchain{}// 시작할떄 한번만 실행 시키고 싶다 Sync 패키지를 사용한다. 어떤 스레드가 있어도 누구든 한번만 수행 sync.Once
		once.Do(func() {
			b = &blockchain{}
			b.blocks = append(b.blocks, createBlock("Genesis Block"))
		})
	}
	return b
}

func getLastHash() string {
	// if len(b.blocks) > 0 {// singleton에서 instance는 getblockchain으로만 가져오다
	if len(GetBlockChain().blocks) != 0 {
		return GetBlockChain().blocks[len(GetBlockChain().blocks)-1].hash
	}
	return ""
}
func (b *blockchain) AddBlock(data string) {
	//hash 가져와야함
	newBlock := block{data, "", b.getLastHash()}
	//hash 생성해야함
	hash := sha256.Sum256([]byte(newBlock.data + newBlock.hash))
	newBlock.hash = fmt.Sprintf("%x", hash)
	//Block에 추가해야함
	b.blocks = append(b.blocks, newBlock)
}
func (b *blockchain) ListBlocks() {
	for _, block := range b.blocks {
		fmt.Printf("Data:%s ", block.data)
		fmt.Printf("Hash:%s ", block.hash)
		fmt.Printf("Previous Hash:%s \n", block.previousHash)

	}
}

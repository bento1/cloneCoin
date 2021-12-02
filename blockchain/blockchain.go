package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

//ver1 21-11-25 only blockchain develop
//여러 기능을 추가하면서 refactoring을 계속 할 것임 ... transaction, database..
type Block struct {
	Data         string `json:"data"` //transaction 등이 바뀔수 있다.
	Hash         string `json:"hash"`
	PreviousHash string `json:"previoushash,omitempty"`
	Height       int    `json:"height"`
}
type blockchain struct {
	// blocks []block
	blocks []*Block // 복사하고싶지 않음
	// Block의 array로 저장하고 있는데, 이는 메모리에 올라가 있음
	// 재시작되면 Genesis부터 시작하게된다.
	// Blockchain 은 가장 마지막의 hash의 포인터만 남긴다.
	// 븍록체인 구조를 DB에 검색할 수 있게하는 메소드 집합체로만드는것이 더좋을것
	// 예를 들면 DB에서 특정 해쉬를 찾는 함수를 만드는거나
	// 아님 블록전체를 가져오는 것을 만들거나
	// 모든 블럭은 url height로 저장됨
}

var b *blockchain //singleton pattern 이녀석을  외부에서 읽게함-> 1개의 인스턴스만 존재하게됨
var once sync.Once

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockChain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}
func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PreviousHash))
	b.Hash = fmt.Sprintf("%x", hash)
}
func GetBlockChain() *blockchain { //singleton pattern 인스턴스를 외부 읽는 메소드-> 1개의 인스턴스만 존재하게됨
	if b == nil {
		// b = &blockchain{}// 시작할떄 한번만 실행 시키고 싶다 Sync 패키지를 사용한다. 어떤 스레드가 있어도 누구든 한번만 수행 sync.Once
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}

func getLastHash() string {
	// if len(b.blocks) > 0 {// singleton에서 instance는 getblockchain으로만 가져오다
	if len(GetBlockChain().blocks) != 0 {
		return GetBlockChain().blocks[len(GetBlockChain().blocks)-1].Hash
	}
	return ""
}

func (b *blockchain) AddBlock(data string) {
	//hash 가져와야함
	// newBlock := block{data, "", getLastHash()}
	//hash 생성해야함
	// hash := sha256.Sum256([]byte(newBlock.data + newBlock.hash))
	// newBlock.hash = fmt.Sprintf("%x", hash)
	//Block에 추가해야함
	b.blocks = append(b.blocks, createBlock(data))
}
func (b *blockchain) ListBlocks() []*Block {
	return GetBlockChain().blocks
}

var ErrNotFound = errors.New("block nor found")

func (b *blockchain) GetBlock(height int) (*Block, error) {
	if height > len(b.blocks) {
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}

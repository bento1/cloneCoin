package main

import (
	"github.com/cloneCoin/blockchain"
)

// B1
// 	b1Hash=(data+"")
// B2
// 	b2Hash=(data2+b1Hash)
//B1 의 내용이 바뀌면 B2 해쉬가 전ㄴ체가 다 다름 무효
// 위에처럼 블록을 링크 할수 있음
func main() {
	// genesisBlock := block{"Genesis Block", (""), ("")} //genesis block means Init Block of Network
	// // genesisBlock.hash=fn(genesisBlock.data+genesisBlock.previousHash) 해쉬를 생성하는 방법? 이라고 임의로 정의함
	// // genesisBlock.hash=sha256.Sum256(genesisBlock.data+genesisBlock.previousHash) [] byte를 인풋으로 받음
	// hash_ := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.previousHash))
	// // fmt.Println(hash_) 32 길이의 BYTE ARRAY 리턴됨
	// // fmt.Printf("%x", hash_) // 비트코인 이더리움 해쉬는 Hexadecimal, 16진수로 되어있음
	// hash := fmt.Sprintf("%x", hash_) //hexaHash로 변경
	// genesisBlock.hash = hash
	// fmt.Println(genesisBlock)
	chain := blockchain.GetBlockChain()
	// chain.AddBlock("Genesis Block")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	// chain.ListBlocks()
	// for _, block := range presentation {
	// 	fmt.Println(block.data) //block내  멤버를 export해야함.  불합리하다고 생각함 ㅎㅎ
	// }
	//http expolorer, json api, cli (command line interface) 추가 구현 예정
}

// func getBytea(word string) []byte{
// 	var result [] byte ;
// 	for _,value := range word{
// 		result = append(result, value)
// 	}
// 	return result[]
// }

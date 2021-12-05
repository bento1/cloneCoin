package main

import "github.com/bento1/cloneCoin/blockchain"

// "github.com/bento1/cloneCoin/cli"

func main() {
	// cli.Start()
	blockchain.BlockChain().AddBlock("First Block")
	blockchain.BlockChain().AddBlock("Second Block")
	blockchain.BlockChain().AddBlock("Third Block")
}

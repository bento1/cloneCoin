package main

import (
	"github.com/github.com/bento1/cloneCoin/blockchain"
	"github.com/github.com/bento1/cloneCoin/cli"
	"github.com/github.com/bento1/cloneCoin/db"
)

func main() {
	defer db.Close()
	blockchain.BlockChain()
	cli.Start()

}

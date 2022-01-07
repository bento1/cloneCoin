package main

import (
	"github.com/github.com/bento1/cloneCoin/cli"
	"github.com/github.com/bento1/cloneCoin/db"
)

func main() {
	defer db.Close()
	cli.Start()

}

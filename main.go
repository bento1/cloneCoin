package main

import (
	"time"

	"github.com/bento1/cloneCoin/explorer"

	"github.com/bento1/cloneCoin/rest"
)

func main() {
	go explorer.Start(3333) //go를 넣어야 프로세스가 2개로 실행됨.
	go rest.Start(4000)

	time.Sleep(time.Second * 10)
}

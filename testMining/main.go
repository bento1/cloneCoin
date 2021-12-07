package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func main() {
	difficulty := 2
	target := strings.Repeat("0", difficulty) // strings.Repeat는 원하는 스트링을 반복하게함

	Nonce := 0
	for {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte("hello"+fmt.Sprint(Nonce))))
		fmt.Println("target : ", target)
		fmt.Println("nonce : ", Nonce)
		fmt.Println("hash : ", hash)
		if strings.HasPrefix(hash, target) {
			return
		}
		Nonce = Nonce + 1
	}

}

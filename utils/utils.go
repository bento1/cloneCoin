package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytea(i interface{}) []byte {
	var blockbuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockbuffer)
	HandleErr(encoder.Encode(i))
	return blockbuffer.Bytes()
}

func FromBytea(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i)) //없으면 nil return함
}

func Hash(i interface{}) string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v", i)))) //v는 인터페이스 formatting해줌
	return hash
}

func Splitter(s string, sep string, i int) string {
	r := strings.Split(s, sep)
	if i > len(r)-1 {
		return ""
	}
	return r[i]
}

func ToJson(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}

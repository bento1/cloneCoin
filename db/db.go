package db

import (
	"fmt"

	"github.com/bento1/cloneCoin/utils"
	"github.com/boltdb/bolt"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

// var once sync.Once

func DB() *bolt.DB {
	if db == nil {
		// once.Do()
		dbPointer, err := bolt.Open(dbName, 0600, nil)

		utils.HandleErr(err)

		db = dbPointer
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))

			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func SaveBlock(key string, value []byte) { //key는 hash가 됨, value는 block을 저장함
	fmt.Printf("Saving Block\nhash: %s\ndata %b\n", key, value)
	err := DB().Update(func(t *bolt.Tx) error {
		//bucker에 저장함
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(key), value)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockChain(data []byte) { //마지막 해쉬와 Height가 담긴 blockchain이 저장된다.
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func BlockChain() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

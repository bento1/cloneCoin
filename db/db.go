package db

import (
	"github.com/bento1/cloneCoin/utils"

	bolt "github.com/boltdb/bolt"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
)

var db *bolt.DB

// var once sync.Once

func DB() *bolt.DB {
	if db == nil {
		// once.Do()
		dbPointer, err := bolt.Open("blockchain.db", 0600, nil)
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

package bbolt_store

import (
	"sync"

	"github.com/iAziz786/raft/config"
	bolt "go.etcd.io/bbolt"
)

type Store struct {
	db *bolt.DB
}

var once sync.Once
var db *bolt.DB

const bucketName = "store"

func NewStore() *Store {
	var dbName = config.Name + ".db"
	once.Do(func() {
		tempDB, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			panic(err)
		}
		db = tempDB
	})

	return &Store{
		db: db,
	}
}

func (s *Store) Get(key string) (val []byte, err error) {
	s.db.View(func(txn *bolt.Tx) error {
		b := txn.Bucket([]byte(bucketName))
		val = b.Get([]byte(key))

		return nil
	})
	return
}

func (s *Store) Set(key string, val []byte) (err error) {
	s.db.Update(func(txn *bolt.Tx) error {
		b, er := txn.CreateBucketIfNotExists([]byte(bucketName))
		if er != nil {
			err = er
			return er
		}
		er = b.Put([]byte(key), val)
		if er != nil {
			err = er
			return er
		}

		return nil
	})
	return
}

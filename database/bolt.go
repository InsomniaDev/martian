package database

import (
	"encoding/json"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

var MartianData Database

type Database struct {
	Connection *bolt.DB
}

var (
	IntegrationBucket  = []byte("integration_bucket")
	SubscriptionBucket = []byte("subscription_bucket")
	DayMemoryBucket    = []byte("day_memory_bucket")
)

func init() {
	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	MartianData.Connection = db
}

// func (d *Database) Init() {
// 	db, err := bolt.Open("./config/martian.db", 0600, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	d.Connection = db
// 	defer d.Connection.Close()
// }

// GetSubscriptionValues will retrieve all subscription values in the provided bucket
func (d *Database) GetSubscriptionValues(key string) (value []byte, err error) {
	return d.getBucketValue(SubscriptionBucket, key)
}

// PutSubscriptionValue will retrieve all subscription values in the provided bucket
func (d *Database) PutSubscriptionValue(key string, value []byte) (err error) {
	return d.putBucketValue(SubscriptionBucket, key, value)
}

// putBucketValue will insert a new integration value into the database
func (d *Database) putBucketValue(bucket []byte, key string, value interface{}) error {
	// db, err := bolt.Open("./config/martian.db", 0600, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	err := d.Connection.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			log.Fatal(err)
			return err
		}

		byteValue, err := json.Marshal(value)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = bucket.Put([]byte(key), byteValue)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// getBucketValue will retrieve the respective integration value from the database
func (d *Database) getBucketValue(bucket []byte, key string) (value []byte, err error) {
	err = d.Connection.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			fmt.Println(err)
			return err
		}

		value = bucket.Get([]byte(key))
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return value, nil
}

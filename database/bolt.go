package database

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

type Database struct {
	Connection *bolt.DB
}

var (
	IntegrationBucket  = []byte("integration_bucket")
	SubscriptionBucket = []byte("subscription_bucket")
)

func (d *Database) Init() {
	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	d.Connection = db
	defer d.Connection.Close()
}

// GetSubscriptionValues will retrieve all subscription values in the provided bucket
func (d *Database) GetSubscriptionValues(key string) (value interface{}, err error) {
	return d.getBucketValue(SubscriptionBucket, key)
}

// getBucketValue will retrieve the respective integration value from the database
func (d *Database) getBucketValue(bucket []byte, key string) (value interface{}, err error) {
	err = d.Connection.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			log.Fatal(err)
			return err
		}

		value = bucket.Get([]byte(key))
		return nil
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return value, nil
}

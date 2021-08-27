package database

import (
	"encoding/json"
	"log"

	bolt "go.etcd.io/bbolt"
)

type Database struct {
	Connection *bolt.DB
}

var (
	IntegrationBucket = []byte("integration_bucket")
)

func (d *Database) Init() {
	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	d.Connection = db
	defer d.Connection.Close()
}

// RetrieveAllValuesInBucket will retrieve all of the values in the provided bucket
func (d *Database) RetrieveAllValuesInBucket(bucket []byte) (value map[string]string, err error) {
	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	foundIntegrations := make(map[string]string)
	err = db.View(func(tx *bolt.Tx) error {
		iter := tx.Bucket(bucket)
		if iter == nil {
			return err
		}

		c := iter.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			foundIntegrations[string(k)] = string(v)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	value = foundIntegrations
	return value, nil
}

// PutIntegrationValue will insert a new integration value into the database
func (d *Database) PutIntegrationValue(key string, value interface{}) error {
	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(IntegrationBucket)
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

// GetIntegrationValue will retrieve the respective integration value from the database
func (d *Database) GetIntegrationValue(key string) (value interface{}, err error) {

	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(IntegrationBucket)
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

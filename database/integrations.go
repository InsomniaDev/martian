package database

import (
	"encoding/json"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

// RetrieveAllValuesInBucket will retrieve all of the values in the provided bucket
func (d *Database) RetrieveAllValuesInBucket(bucket []byte) (value map[string]string, err error) {
	// db, err := bolt.Open("./config/martian.db", 0600, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	foundIntegrations := make(map[string]string)
	err = d.Connection.View(func(tx *bolt.Tx) error {
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
	// db, err := bolt.Open("./config/martian.db", 0600, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	err := d.Connection.Update(func(tx *bolt.Tx) error {
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

// DeleteIntegrationValue will retrieve the respective integration value from the database
func (d *Database) DeleteIntegrationValue(key string) (err error) {

	// db, err := bolt.Open("./config/martian.db", 0600, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	d.Connection.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(IntegrationBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

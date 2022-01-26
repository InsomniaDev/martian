package database

import (
	"encoding/json"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

// GetDayMemoryByHour will get the memory of the day by hour
func (d *Database) GetDayMemoryByHour(key string) (value []byte, err error) {
	return d.getBucketValue(DayMemoryBucket, key)
}

// DeleteMemoryHourFromDay will delete the memory from the bucket for today
func (d *Database) DeleteMemoryHourFromDay(key string) (err error) {

	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(DayMemoryBucket)
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

// RetrieveAllMemories will return all of the memories that are in the bucket
// CURRENTLY NOT USED
func (d *Database) RetrieveAllMemories() (value map[string][]byte, err error) {
	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	foundIntegrations := make(map[string][]byte)
	err = db.View(func(tx *bolt.Tx) error {
		iter := tx.Bucket(DayMemoryBucket)
		if iter == nil {
			return err
		}

		c := iter.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			foundIntegrations[string(k)] = v
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

// InsertMemory will insert a new a new memory
func (d *Database) InsertMemory(key string, value interface{}) error {
	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(DayMemoryBucket)
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

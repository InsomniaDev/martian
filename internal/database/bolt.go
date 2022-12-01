package database

import (
	"encoding/json"
	"errors"

	bolt "github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

var MartianData Database

type Database struct {
	Connection *bolt.DB
}

var (
	MemoryBucket      = []byte("memory")
	DeviceBucket      = []byte("device")
	DeviceGraphBucket = []byte("deviceGraph")
	TimeTableBucket   = []byte("timeChain")
)

func init() {
	db, err := bolt.Open("./config/martian.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	MartianData.Connection = db
}

// GetMemoryData will grab the memory data from the database
// 		returns a boolean on if the value exists along with the value
func (d *Database) GetMemoryData() (bool, []byte) {
	switch value, err := d.getBucketValue(MemoryBucket, "memory"); {
	case err != nil:
		return false, nil
	default:
		return true, value
	}
}

// StoreMemoryData will store the memory data in the bolt database
//		returns a boolean on storage success
func (d *Database) StoreMemoryData(jsonBrain interface{}) bool {
	if err := d.putBucketValue(MemoryBucket, "memory", jsonBrain); err != nil {
		return false
	}
	return true
}

func (d *Database) RecreateGraphBucket() error {
	if err := d.recreateBucket(DeviceGraphBucket); err != nil {
		return err
	}
	return nil
}

func (d *Database) RecreateTimeBucket() error {
	if err := d.recreateBucket(TimeTableBucket); err != nil {
		return err
	}
	return nil
}

// GetDeviceValues
func (d *Database) GetDeviceValues(uniqueHash string) (device Device, err error) {
	value, err := d.getBucketValue(DeviceBucket, uniqueHash)
	if err != nil {
		return Device{}, err
	}
	json.Unmarshal(value, &device)
	return
}

// PutDeviceValues
func (d *Database) PutDeviceValues(device Device) (err error) {
	return d.putBucketValue(DeviceBucket, device.UniqueHash, device)
}

// GetAllDevices
func (d *Database) GetAllDevices() (devices []Device) {
	entries := d.getAllKeys(string(DeviceBucket))
	for _, entry := range entries {
		var device Device
		err := json.Unmarshal(entry, &device)
		if err != nil {
			log.Fatal(err)
		}
		devices = append(devices, device)
	}
	return
}

// GetDeviceGraphValues
func (d *Database) GetDeviceGraphValues(fromUniqueHash string) (deviceGraphs UniqueHashGraphs, err error) {
	value, err := d.getBucketValue(DeviceGraphBucket, fromUniqueHash)
	if err != nil {
		return UniqueHashGraphs{}, err
	}
	json.Unmarshal(value, &deviceGraphs)
	return
}

// PutDeviceGraphValues
func (d *Database) PutDeviceGraphValues(graphs UniqueHashGraphs) (err error) {
	if len(graphs.Graphs) == 0 {
		return errors.New("Requires graphs to be populated")
	}
	return d.putBucketValue(DeviceGraphBucket, graphs.Graphs[0].FromUniqueHash, graphs)
}

// GetAllDeviceGraphss
func (d *Database) GetAllDeviceGraphs() (graphs []DeviceGraph) {
	entries := d.getAllKeys(string(DeviceGraphBucket))
	for _, entry := range entries {
		var graph UniqueHashGraphs
		err := json.Unmarshal(entry, &graph)
		if err != nil {
			log.Fatal(err)
		}
		for i := range graph.Graphs {
			graphs = append(graphs, graph.Graphs[i])
		}
	}
	return
}

// GetAllDeviceGraphss
func (d *Database) GetAllGraphs() (graphs []UniqueHashGraphs) {
	entries := d.getAllKeys(string(DeviceGraphBucket))
	for _, entry := range entries {
		var graph UniqueHashGraphs
		err := json.Unmarshal(entry, &graph)
		if err != nil {
			log.Fatal(err)
		}
		graphs = append(graphs, graph)
	}
	return
}

// GetAllTimeBuckets
func (d *Database) getAllTimeBuckets() (timeTables []TimeTable) {
	entries := d.getAllKeys(string(TimeTableBucket))
	for _, entry := range entries {
		var timeTable TimeTable
		err := json.Unmarshal(entry, &timeTable)
		if err != nil {
			log.Fatal(err)
		}
		timeTables = append(timeTables, timeTable)
	}
	return
}

// GetTimeTableValues
func (d *Database) GetTimeTableValues(timeBlock string) (times TimeTable, err error) {
	value, err := d.getBucketValue(TimeTableBucket, timeBlock)
	if err != nil {
		return TimeTable{}, err
	}
	json.Unmarshal(value, &times)
	return
}

// PutTimeTableValues
func (d *Database) PutTimeTableValues(timeTable string, timeEntries TimeTable) (err error) {
	if len(timeEntries.Times) == 0 {
		return errors.New("requires times to be populated")
	}
	return d.putBucketValue(TimeTableBucket, timeTable, timeEntries)
}

// putBucketValue will insert a new integration value into the database
func (d *Database) putBucketValue(bucket []byte, key string, value interface{}) error {
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
	// d.SqlConnectionPool.Acquire()
	err = d.Connection.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			log.Error("get bucket err:", bucket, err)
			return err
		}

		value = bucket.Get([]byte(key))
		return nil
	})

	if err != nil {
		log.Error("get bucket err:", bucket, err)
		return nil, err
	}
	return value, nil
}

// recreateBucket will delete everything from the provided bucket and create a new blank bucket
func (d *Database) recreateBucket(bucketKey []byte) error {

	err := d.Connection.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketKey)
		if err != nil {
			log.Fatal(err)
			return err
		}
		_, err = tx.CreateBucketIfNotExists(bucketKey)
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

func (d *Database) getAllKeys(bucket string) (entries [][]byte) {
	err := d.Connection.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Fatal(err)
			return err
		}

		bucket.ForEach(func(k, v []byte) error {
			entries = append(entries, v)
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return
}

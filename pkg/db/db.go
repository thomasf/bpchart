package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/thomasf/bpchart/pkg/omron"
	"github.com/thomasf/lg"
)

// EntryBucket .
type DB struct {
	*bolt.DB
	BucketName []byte
}

func entryKey(entry omron.Entry) []byte {
	return []byte(entry.Time.UTC().Format(time.RFC3339))
}

func (db *DB) All() ([]omron.Entry, error) {
	var allEntries []omron.Entry
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.BucketName)

		err := b.ForEach(func(k []byte, v []byte) error {
			var entry omron.Entry
			err := json.Unmarshal(v, &entry)
			if err != nil {
				return err
			}
			allEntries = append(allEntries, entry)
			return nil
		})
		return err
	})
	return allEntries, err
}

func (db *DB) SaveEntry(entry omron.Entry, b *bolt.Bucket) error {
	encoded, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return b.Put([]byte(entryKey(entry)), encoded)
}

func (db *DB) SaveEntries(all []omron.Entry) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(db.BucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		for _, v := range all {
			key := entryKey(v)
			data := b.Get(key)
			if data == nil {
				err := db.SaveEntry(v, b)
				if err != nil {
					return err
				}

			} else {
				lg.V(10).Infoln(key, "already exist")
			}

		}
		return nil
	})

}

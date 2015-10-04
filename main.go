package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sort"
	"time"

	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/thomasf/bpchart/pkg/omron"
	"github.com/thomasf/lg"
)

// devleopment mode
const FAKE = true

var entriesBucketName = "entries"

func init() {
	if FAKE {
		entriesBucketName = "fakeEntries"
	}

}

func entryKey(entry omron.Entry) []byte {
	return []byte(entry.Time.UTC().Format(time.RFC3339))
}
func saveEntry(entry omron.Entry, b *bolt.Bucket) error {
	encoded, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return b.Put([]byte(entryKey(entry)), encoded)
}

// EntryBucket .
type EntryBucket struct {
	*bolt.Bucket
}

func fakeImportFromDevice(db *bolt.DB) error {

	genRandomEntry := func(lastTime time.Time) omron.Entry {
		return omron.Entry{
			Time:  lastTime.Add(8 + time.Hour*12),
			Sys:   100 + rand.Intn(60),
			Dia:   60 + rand.Intn(35),
			Pulse: 60 + rand.Intn(40),
			Bank:  rand.Intn(1),
		}
	}

	genIncrEntry := func(e omron.Entry) omron.Entry {
		k := rand.Intn(10) + 1
		return omron.Entry{
			Time:  e.Time.Add((15 + time.Duration(rand.Intn(60))) * time.Second),
			Sys:   e.Sys - 5 + rand.Intn(k),
			Dia:   e.Dia - 5 + rand.Intn(k),
			Pulse: e.Pulse - 5 + rand.Intn(k),
			Bank:  e.Bank,
		}

	}

	all := make([]omron.Entry, 50)
	t := time.Now().Add(-300 * time.Hour * 24)
	for i := range all {
		if i%3 == 0 {
			e := genRandomEntry(t)
			t = e.Time
			all[i] = e
		} else {
			e := genIncrEntry(all[i-1])
			t = e.Time
			all[i] = e
		}

	}

	sort.Sort(omron.ByTime(all))
	// data, err := json.MarshalIndent(all, "", "  ")
	// lg.Infoln(string(data))
	// if err != nil {
	// lg.Fatal(err)
	// }
	// os.Exit(1)

	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(entriesBucketName))
		return nil
	})

	return saveEntries(db, all)
}

func importFromDevice(db *bolt.DB) error {
	if err := omron.Open(); err != nil {
		return err
	}
	defer func() {
		if err := omron.Close(); err != nil {
			lg.Errorln(err)
		}
	}()

	var all []omron.Entry

	{
		entries, err := omron.Read(0)
		if err != nil {
			return err
		}
		all = append(all, entries...)

	}

	{
		entries, err := omron.Read(1)
		if err != nil {
			return err
		}
		all = append(all, entries...)
	}
	sort.Sort(omron.ByTime(all))
	// data, err := json.MarshalIndent(all, "", "  ")
	// if err != nil {
	// lg.Fatal(err)
	// }
	return saveEntries(db, all)

}

func saveEntries(db *bolt.DB, all []omron.Entry) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(entriesBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		for _, v := range all {
			key := entryKey(v)
			data := b.Get(key)
			if data == nil {
				err := saveEntry(v, b)
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

func httpServer(db *bolt.DB) error {

	http.HandleFunc("/json/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var entries []omron.Entry
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(entriesBucketName))

			b.ForEach(func(k []byte, v []byte) error {
				var entry omron.Entry
				err := json.Unmarshal(v, &entry)
				if err != nil {
					return err
				}
				entries = append(entries, entry)
				return nil
			})
			return nil
		})
		data, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			lg.Fatal(err)
		}
		w.Write(data)
	})

	http.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(
				http.Dir("build/assets"))))

	http.Handle("/",
		http.FileServer(
			http.Dir("browser/html")))

	return http.ListenAndServe(":8080", nil)

}
func main() {
	runtime.LockOSThread()
	flag.Set("logtostderr", "true")
	lg.CopyStandardLogTo("INFO")
	lg.SetSrcHighlight("libomron")
	flag.Parse()

	db, err := bolt.Open("bpchart.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		lg.Fatal(err)
	}
	if FAKE {
		err = fakeImportFromDevice(db)
	} else {
		err = importFromDevice(db)
	}
	if err != nil {
		lg.Fatal(err)
	}

	lg.Fatal(httpServer(db))

}

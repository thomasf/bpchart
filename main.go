package main

import (
	"encoding/json"
	"flag"

	"math/rand"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/thomasf/bpchart/pkg/db"
	"github.com/thomasf/bpchart/pkg/omron"
	omronread "github.com/thomasf/bpchart/pkg/readC"
	"github.com/thomasf/bpchart/pkg/score"
	"github.com/thomasf/lg"
)

// devleopment mode
const FAKE = false

var entriesBucketName = "entries"

func init() {
	if FAKE {
		entriesBucketName = "fakeEntries"
	}
}

func fakeImportFromDevice(db *db.DB) error {

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

	all := make([]omron.Entry, 500)
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

	return db.SaveEntries(all)
}

func importFromDevice(db *db.DB) error {
	if err := omronread.Open(); err != nil {
		return err
	}
	defer func() {
		if err := omronread.Close(); err != nil {
			lg.Errorln(err)
		}
	}()

	var all []omron.Entry

	{
		entries, err := omronread.Read(0)
		if err != nil {
			return err
		}
		all = append(all, entries...)

	}

	{
		entries, err := omronread.Read(1)
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
	return db.SaveEntries(all)

}

func httpServer(db *db.DB) error {

	http.HandleFunc("/json/", func(w http.ResponseWriter, r *http.Request) {

		dtMin := time.Now().Add(-time.Hour * 24 * 7 * 4)
		dtMax := time.Now().Add(time.Minute)
		{
			type stp struct {
				t  *time.Time
				qp string
			}
			for _, v := range []stp{
				{&dtMin, "dt_min"},
				{&dtMax, "dt_max"},
			} {
				ts := r.URL.Query().Get(v.qp)
				t, err := time.Parse(time.RFC3339, ts)
				if err != nil {
					t, err = time.Parse("2006-01-02", ts)
					if err != nil {
						continue
					}
				}
				*v.t = t
			}
		}

		avgMinutes := 10
		if r.URL.Query().Get("avg_minutes") != "" {
			var err error
			avgMinutes, err = strconv.Atoi(r.URL.Query().Get("avg_minutes"))
			if err != nil {
				lg.Fatalln(err)
			}
		}

		allEntries, err := db.All()

		var filteredEntires []omron.Entry
		for _, entry := range allEntries {
			if entry.Time.After(dtMin) && entry.Time.Before(dtMax) {
				filteredEntires = append(filteredEntires, entry)
			}
		}

		avgEntries := omron.AvgWithinDuration(
			filteredEntires,
			time.Duration(avgMinutes)*time.Minute)

		scoredEntries := score.All(avgEntries)

		w.Header().Set("Content-type", "application/json")
		data, err := json.MarshalIndent(scoredEntries, "", "  ")
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

	bdb, err := bolt.Open("bpchart.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		lg.Fatal(err)
	}
	db := &db.DB{DB: bdb, BucketName: []byte(entriesBucketName)}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		lg.Fatal(httpServer(db))
		wg.Done()
	}()

	wg.Add(1)
	go func() {

		if FAKE {
			err = fakeImportFromDevice(db)
		} else {
			err = importFromDevice(db)
		}
		if err != nil {
			lg.Fatal(err)
		}
		wg.Done()
	}()

	wg.Wait()
}

package omron

import "time"

type ByTime []Entry

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }

// Entry .
type Entry struct {
	Time  time.Time `json:"time"`  // recorded time
	Sys   int       `json:"sys"`   // mmHg
	Dia   int       `json:"dia"`   // mmHg
	Pulse int       `json:"pulse"` // beats per minute
	Bank  int       `json:"bank"`  // bank
}

func AvgWithinDuration(entries []Entry, duration time.Duration) []Entry {
	var avgEntries []Entry

	var prev Entry
	var sum Entry
	nSum := 0

	for _, v := range entries {
		if v.Time.Before(prev.Time.Add(duration)) {
			if nSum == 0 {
				sum.Dia = prev.Dia
				sum.Sys = prev.Sys
				sum.Pulse = prev.Pulse
			}
			sum.Time = v.Time
			sum.Dia += v.Dia
			sum.Sys += v.Sys
			sum.Pulse += v.Pulse
			nSum++
		} else if nSum != 0 {
			e := Entry{
				Dia:   int(float64(sum.Dia) / float64(nSum+1)),
				Sys:   int(float64(sum.Sys) / float64(nSum+1)),
				Pulse: int(float64(sum.Pulse) / float64(nSum+1)),
				Time:  sum.Time,
			}
			avgEntries = append(avgEntries, e)
			nSum = 0
		} else {
			avgEntries = append(avgEntries, v)
		}
		prev = v

	}

	return avgEntries
}

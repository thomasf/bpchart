package score

import "github.com/thomasf/bpchart/pkg/omron"

//go:generate stringer -type=Score
type Score int

const (
	OK Score = 1 + iota
	LowWarning
	Low
	HighWarning
	High
)

// Entry is a scored entry.
type Entry struct {
	omron.Entry
	SysScore   Score `json:"sysScore"`
	DiaScore   Score `json:"diaScore"`
	PulseScore Score `json:"pulseScore"`
}

func All(entries []omron.Entry) []Entry {
	var scoredEntries []Entry
	for _, v := range entries {
		scoredEntries = append(scoredEntries, New(v))

	}
	return scoredEntries
}

func New(e omron.Entry) Entry {
	var SysScore, DiaScore, PulseScore Score
	switch {
	case e.Sys > 139:
		SysScore = High
	case e.Sys > 119:
		SysScore = HighWarning
	case e.Sys < 90:
		SysScore = LowWarning
	default:
		SysScore = OK
	}

	switch {
	case e.Dia > 89:
		DiaScore = High
	case e.Dia > 79:
		DiaScore = HighWarning
	case e.Dia < 60:
		DiaScore = LowWarning
	default:
		DiaScore = OK
	}

	switch {
	case e.Pulse > 90:
		PulseScore = High
	case e.Pulse > 80:
		PulseScore = HighWarning
	default:
		PulseScore = OK
	}

	return Entry{
		Entry:      e,
		SysScore:   SysScore,
		DiaScore:   DiaScore,
		PulseScore: PulseScore,
	}
}

func (s *Score) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

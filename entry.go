package zltest

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

// KeyStatus represents a status of searching and deserialization
// of a key in log entry.
type KeyStatus string

const (
	// KeyFound is used when Key found successfully.
	KeyFound KeyStatus = "KeyFound"
	// KeyBadType is used when Key found, but it's not of expected type.
	KeyBadType KeyStatus = "KeyBadType"
	// KeyMissing is used when Key is not in the log entry.
	KeyMissing KeyStatus = "KeyMissing"
	// KeyBadFormat is used when Key found but its format is wrong.
	KeyBadFormat KeyStatus = "KeyBadFormat"
)

// Entry represents zerolog log entry.
type Entry struct {
	raw string                 // Entry as it was written to the writer.
	m   map[string]interface{} // JSON decoded log entry.
	t   T                      // Test manager.
}

// String implements fmt.Stringer interface and returns log entry
// as it was written to the writer.
func (ent *Entry) String() string {
	return ent.raw
}

// ExpKey tests log entry has a key.
func (ent *Entry) ExpKey(key string) {
	if _, ok := ent.m[key]; !ok {
		ent.t.Errorf("expected %s field to be present", key)
	}
}

// NotExpKey tests log entry has no key.
func (ent *Entry) NotExpKey(key string) {
	if _, ok := ent.m[key]; ok {
		ent.t.Errorf("expected %s field to be not present", key)
	}
}

// Str returns log entry key as a string.
func (ent *Entry) Str(key string) (string, KeyStatus) {
	if itf, ok := ent.m[key]; ok {
		if got, ok := itf.(string); ok {
			return got, KeyFound
		}
		return "", KeyBadType
	}
	return "", KeyMissing
}

// Float64 returns log entry key as a float64.
func (ent *Entry) Float64(key string) (float64, KeyStatus) {
	if itf, ok := ent.m[key]; ok {
		if got, ok := itf.(float64); ok {
			return got, KeyFound
		}
		return 0, KeyBadType
	}
	return 0, KeyMissing
}

// Bool returns log entry key as a boolean.
func (ent *Entry) Bool(key string) (bool, KeyStatus) {
	if itf, ok := ent.m[key]; ok {
		if got, ok := itf.(bool); ok {
			return got, KeyFound
		}
		return false, KeyBadType
	}
	return false, KeyMissing
}

// Time returns log entry key as a time.Time. It uses zerolog.TimeFieldFormat
// to parse the time string representation.
func (ent *Entry) Time(key string) (time.Time, KeyStatus) {
	if itf, ok := ent.m[key]; ok {
		if got, ok := itf.(string); ok {
			tim, err := time.Parse(zerolog.TimeFieldFormat, got)
			if err != nil {
				return time.Time{}, KeyBadFormat
			}
			return tim, KeyFound
		}
		return time.Time{}, KeyBadType
	}
	return time.Time{}, KeyMissing
}

// ExpStr tests log entry has key, its value is a string, and it's equal to exp.
func (ent *Entry) ExpStr(key string, exp string) {
	ent.t.Helper()
	if err := ent.expStr(key, exp); err != "" {
		ent.t.Error(err)
	}
}

func (ent *Entry) expStr(key string, exp string) string {
	ent.t.Helper()
	got, status := ent.Str(key)
	if status == KeyFound {
		if got != exp {
			return fmt.Sprintf(
				"expected entry key '%s' to have value '%s' but got '%s'",
				key,
				exp,
				got,
			)
		}
		return ""
	}
	return formatError(status, key, "string")
}

// ExpTime tests log entry has key, its value is a string representing time in
// zerolog.TimeFieldFormat and it's equal to exp.
func (ent *Entry) ExpTime(key string, exp time.Time) {
	ent.t.Helper()
	if err := ent.expTime(key, exp); err != "" {
		ent.t.Error(err)
	}
}
func (ent *Entry) expTime(key string, exp time.Time) string {
	ent.t.Helper()
	got, status := ent.Time(key)
	if status == KeyFound {
		if !exp.Equal(got) {
			return fmt.Sprintf("expected entry '%s' to be '%s' but is '%s'",
				key,
				exp.Format(zerolog.TimeFieldFormat),
				got.Format(zerolog.TimeFieldFormat),
			)
		}
		return ""
	}
	return formatError(status, key, "string")

}

// ExpDur tests log entry has key and its value is equal to exp time.Duration.
// The duration vale in the entry is multiplied by zerolog.DurationFieldUnit
// before the comparison.
func (ent *Entry) ExpDur(key string, exp time.Duration) {
	ent.t.Helper()
	if err := ent.expDur(key, exp); err != "" {
		ent.t.Error(err)
	}
}
func (ent *Entry) expDur(key string, exp time.Duration) string {
	ent.t.Helper()
	got, status := ent.Float64(key)
	if status == KeyFound {
		gotD := time.Duration(int(got)) * zerolog.DurationFieldUnit
		if gotD != exp {
			return fmt.Sprintf(
				"expected entry key '%s' to have value '%d' (%s) but got '%d' (%s)",
				key,
				exp/zerolog.DurationFieldUnit,
				exp.String(),
				gotD/zerolog.DurationFieldUnit,
				gotD.String(),
			)
		}
		return ""
	}
	return formatError(status, key, "number")
}

// ExpBool tests log entry has a key, its value is boolean and equal to exp.
func (ent *Entry) ExpBool(key string, exp bool) {
	ent.t.Helper()
	if err := ent.expBool(key, exp); err != "" {
		ent.t.Error(err)
	}

}

func (ent *Entry) expBool(key string, exp bool) string {
	ent.t.Helper()
	got, status := ent.Bool(key)
	if status == KeyFound {
		if got != exp {
			return fmt.Sprintf(
				"expected entry key '%s' to have value '%v' but got '%v'",
				key,
				exp,
				got,
			)
		}
		return ""
	}
	return formatError(status, key, "bool")
}

// ExpLoggedWithin tests log entry was logged at exp time. The actual time
// may be within +/- diff.
func (ent *Entry) ExpLoggedWithin(exp time.Time, diff time.Duration) {
	ent.ExpTimeWithin(zerolog.TimestampFieldName, exp, diff)
}

// ExpTimeWithin tests log entry has key, its value is a string representing
// time in zerolog.TimeFieldFormat and it's equal to exp time. The actual time
// may be within +/- diff.
func (ent *Entry) ExpTimeWithin(key string, exp time.Time, diff time.Duration) {
	ent.t.Helper()
	got, status := ent.Time(key)
	if status == KeyFound {
		gotD := math.Abs(float64(exp.Sub(got)))
		if gotD > float64(diff) {
			ent.t.Errorf("expected entry '%s' to be within '%s' but difference is '%s'",
				key,
				diff.String(),
				time.Duration(gotD).String(),
			)
		}
		return
	}
	ent.t.Error(formatError(status, key, "string"))
}

// ExpMsg tests log entry message field (zerolog.MessageFieldName) is
// equal to exp.
func (ent *Entry) ExpMsg(exp string) {
	ent.ExpStr(zerolog.MessageFieldName, exp)
}
func (ent *Entry) expMsg(exp string) string {
	return ent.expStr(zerolog.MessageFieldName, exp)
}

// ExpLevel tests log entry level field (zerolog.LevelFieldName) is
// equal to exp.
func (ent *Entry) ExpLevel(exp zerolog.Level) {
	ent.ExpStr(zerolog.LevelFieldName, exp.String())
}
func (ent *Entry) expLevel(exp zerolog.Level) string {
	return ent.expStr(zerolog.LevelFieldName, exp.String())
}

// ExpNum tests log entry has key and its numerical value is equal to exp.
func (ent *Entry) ExpNum(key string, exp float64) {
	ent.t.Helper()
	if err := ent.expNum(key, exp); err != "" {
		ent.t.Error(err)
	}
}
func (ent *Entry) expNum(key string, exp float64) string {
	ent.t.Helper()
	got, status := ent.Float64(key)
	if status == KeyFound {
		if got != exp {
			expS := strconv.FormatFloat(exp, 'f', -1, 64)
			gotS := strconv.FormatFloat(got, 'f', -1, 64)
			return fmt.Sprintf(
				"expected entry key '%s' to have value '%s' but got '%s'",
				key,
				expS,
				gotS,
			)
		}
		return ""
	}
	return formatError(status, key, "number")
}

// formatError formats error message based on status of log entry key search.
func formatError(status KeyStatus, key, typ string) string {
	switch status {
	case KeyMissing:
		return fmt.Sprintf("expected entry to have key '%s'", key)
	case KeyBadType:
		return fmt.Sprintf("expected entry key '%s' to be '%s'", key, typ)
	case KeyBadFormat:
		return fmt.Sprintf("key '%s' in a wrong format", key)
	default:
		return fmt.Sprintf("invalid KeyStatus '%s'", status)
	}
}

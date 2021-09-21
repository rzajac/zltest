package zltest

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// KeyStatus represents a status of searching and deserialization
// of a key in log entry.
type KeyStatus string

const (
	// KeyFound is used when a field key is found successfully.
	KeyFound KeyStatus = "KeyFound"

	// KeyBadType is used when a field key is found, but it's not of expected type.
	KeyBadType KeyStatus = "KeyBadType"

	// KeyMissing is used when a field key is not in the log entry.
	KeyMissing KeyStatus = "KeyMissing"

	// KeyBadFormat is used when a field key is found but its format is wrong.
	KeyBadFormat KeyStatus = "KeyBadFormat"
)

// Entry represents one zerolog log entry.
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

// ExpKey tests log entry has a field key.
func (ent *Entry) ExpKey(key string) {
	ent.t.Helper()
	if _, ok := ent.m[key]; !ok {
		ent.t.Errorf("expected %s field to be present", key)
	}
}

// NotExpKey tests log entry has no field key.
func (ent *Entry) NotExpKey(key string) {
	ent.t.Helper()
	if _, ok := ent.m[key]; ok {
		ent.t.Errorf("expected %s field to be not present", key)
	}
}

// Str returns log entry field key as a string.
func (ent *Entry) Str(key string) (string, KeyStatus) {
	ent.t.Helper()
	if itf, ok := ent.m[key]; ok {
		if got, ok := itf.(string); ok {
			return got, KeyFound
		}
		return "", KeyBadType
	}
	return "", KeyMissing
}

// ExpStr tests log entry has a field key, its value is a string,
// and it's equal to exp.
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
	return formatError(ent.t, status, key, "string")
}

// ExpStrContains tests log entry has a field key, its value is a string,
// and it contains exp.
func (ent *Entry) ExpStrContains(key string, exp string) {
	ent.t.Helper()
	if err := ent.expStrContains(key, exp); err != "" {
		ent.t.Error(err)
	}
}

func (ent *Entry) expStrContains(key string, exp string) string {
	ent.t.Helper()
	got, status := ent.Str(key)
	if status == KeyFound {
		if !strings.Contains(got, exp) {
			return fmt.Sprintf(
				"expected entry key '%s' to contain '%s' but got '%s'",
				key,
				exp,
				got,
			)
		}
		return ""
	}
	return formatError(ent.t, status, key, "string")
}

// Float64 returns log entry field key as a float64 type.
func (ent *Entry) Float64(key string) (float64, KeyStatus) {
	ent.t.Helper()
	if itf, ok := ent.m[key]; ok {
		if got, ok := itf.(float64); ok {
			return got, KeyFound
		}
		return 0, KeyBadType
	}
	return 0, KeyMissing
}

// Bool returns log entry field key as a boolean type.
func (ent *Entry) Bool(key string) (bool, KeyStatus) {
	ent.t.Helper()
	if itf, ok := ent.m[key]; ok {
		if got, ok := itf.(bool); ok {
			return got, KeyFound
		}
		return false, KeyBadType
	}
	return false, KeyMissing
}

// ExpBool tests log entry has a field key, its value is boolean and equal to exp.
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
	return formatError(ent.t, status, key, "bool")
}

// Time returns log entry field  key as a time.Time. It uses
// zerolog.TimeFieldFormat to parse the time string representation.
func (ent *Entry) Time(key string) (time.Time, KeyStatus) {
	ent.t.Helper()
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

// ExpTime tests log entry has a field key, its value is a string representing
// time in zerolog.TimeFieldFormat and it's equal to exp.
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
	return formatError(ent.t, status, key, "string")

}

// ExpTimeWithin tests log entry has a field key, its value is a string
// representing time in zerolog.TimeFieldFormat and it's equal to exp time.
// The actual time may be within +/- diff.
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
	ent.t.Error(formatError(ent.t, status, key, "string"))
}

// ExpDur tests log entry has a field key and its value is equal to exp
// time.Duration. The duration vale in the entry is multiplied by
// zerolog.DurationFieldUnit before the comparison.
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
	return formatError(ent.t, status, key, "number")
}

// ExpLoggedWithin tests log entry was logged at exp time. The actual time
// may be within +/- diff.
func (ent *Entry) ExpLoggedWithin(exp time.Time, diff time.Duration) {
	ent.t.Helper()
	ent.ExpTimeWithin(zerolog.TimestampFieldName, exp, diff)
}

// ExpMsg tests log entry message field (zerolog.MessageFieldName) is
// equal to exp.
func (ent *Entry) ExpMsg(exp string) {
	ent.t.Helper()
	ent.ExpStr(zerolog.MessageFieldName, exp)
}

// ExpError tests log entry message field (zerolog.ErrorFieldName) is
// equal to exp.
func (ent *Entry) ExpError(exp string) {
	ent.t.Helper()
	ent.ExpStr(zerolog.ErrorFieldName, exp)
}

// ExpLevel tests log entry level field (zerolog.LevelFieldName) is
// equal to exp.
func (ent *Entry) ExpLevel(exp zerolog.Level) {
	ent.t.Helper()
	ent.ExpStr(zerolog.LevelFieldName, exp.String())
}

// ExpNum tests log entry has a field key and its numerical value is equal to exp.
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
	return formatError(ent.t, status, key, "number")
}

// Map returns log entry key as a map.
func (ent *Entry) Map(key string) (map[string]interface{}, KeyStatus) {
	ent.t.Helper()
	if itf, ok := ent.m[key]; ok {
		if got, ok := itf.(map[string]interface{}); ok {
			return got, KeyFound
		}
		return nil, KeyBadType
	}
	return nil, KeyMissing
}

// formatError formats error message based on status of log entry key search.
func formatError(t T, status KeyStatus, key, typ string) string {
	t.Helper()
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

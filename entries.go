package zltest

import "time"

// Entries represents collection of zerolog log entries.
type Entries struct {
	e []*Entry // Log entries.
	t T        // Test manager.
}

// Get returns the list of Entry in Entries
func (ets Entries) Get() []*Entry {
	return ets.e
}

// ExpEntry returns nth logged entry.
func (ets Entries) ExpEntry(n int) *Entry {
	ets.t.Helper()
	if n < len(ets.e) {
		return ets.e[n]
	}
	ets.t.Fatalf("expected %d%s logged entry to exist", n, ordinal(n))
	return nil
}

// ordinal returns English ordinal for a whole number.
func ordinal(n int) string {
	switch n {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

// ExpLen tests that there is want number of entries.
func (ets Entries) ExpLen(want int) {
	ets.t.Helper()
	have := len(ets.e)
	if have != want {
		ets.t.Errorf("expected %d entries got %d", want, have)
	}
}

// ExpStr tests that at least one log entry has key, its value is a
// string, and it's equal to exp.
func (ets Entries) ExpStr(key string, exp string) {
	ets.t.Helper()
	ets.exp(func(e *Entry) string { return e.expStr(key, exp) })
}

// ExpStrContains tests that at least one log entry has key, its value is a
// string, and it contains exp.
func (ets Entries) ExpStrContains(key string, exp string) {
	ets.t.Helper()
	ets.exp(func(e *Entry) string { return e.expStrContains(key, exp) })
}

// NotExpStr tests that no log entry has key, its value is a
// string, and it's equal to exp.
func (ets Entries) NotExpStr(key string, exp string) {
	ets.t.Helper()
	ets.notExp(func(e *Entry) string { return e.expStr(key, exp) })
}

// ExpTime tests that at least one log entry has key, its value is a
// string representing time in zerolog.TimeFieldFormat and it's equal
// to exp.
func (ets Entries) ExpTime(key string, exp time.Time) {
	ets.t.Helper()
	ets.exp(func(e *Entry) string { return e.expTime(key, exp) })
}

// NotExpTime tests that no one log entry has key, its value is a
// string representing time in zerolog.TimeFieldFormat and it's equal
// to exp.
func (ets Entries) NotExpTime(key string, exp time.Time) {
	ets.t.Helper()
	ets.notExp(func(e *Entry) string { return e.expTime(key, exp) })
}

// ExpDur tests that at least one log entry has key and its value is
// equal to exp time.Duration.  The duration vale in the entry is
// multiplied by zerolog.DurationFieldUnit before the comparison.
func (ets Entries) ExpDur(key string, exp time.Duration) {
	ets.t.Helper()
	ets.exp(func(e *Entry) string { return e.expDur(key, exp) })
}

// NotExpDur tests that no log entry has key and its value is
// equal to exp time.Duration.  The duration vale in the entry is
// multiplied by zerolog.DurationFieldUnit before the comparison.
func (ets Entries) NotExpDur(key string, exp time.Duration) {
	ets.t.Helper()
	ets.notExp(func(e *Entry) string { return e.expDur(key, exp) })
}

// ExpBool tests that at lest one entry has a key, its value is
// boolean and equal to exp.
func (ets Entries) ExpBool(key string, exp bool) {
	ets.t.Helper()
	ets.exp(func(e *Entry) string { return e.expBool(key, exp) })
}

// NotExpBool tests that no log entry has a key, its value is
// boolean and equal to exp.
func (ets Entries) NotExpBool(key string, exp bool) {
	ets.t.Helper()
	ets.notExp(func(e *Entry) string { return e.expBool(key, exp) })
}

// ExpMsg tests that at least one log entry message field
// (zerolog.MessageFieldName) is equal to exp.
func (ets Entries) ExpMsg(exp string) {
	ets.t.Helper()
	ets.exp(func(e *Entry) string { return e.expMsg(exp) })
}

// NotExpMsg tests that no log entry message field
// (zerolog.MessageFieldName) is equal to exp.
func (ets Entries) NotExpMsg(exp string) {
	ets.t.Helper()
	ets.notExp(func(e *Entry) string { return e.expMsg(exp) })
}

// ExpError tests that at least one log entry error field
// (zerolog.ErrorFieldName) is equal to exp.
func (ets Entries) ExpError(exp string) {
	ets.t.Helper()
	ets.exp(func(e *Entry) string { return e.expError(exp) })
}

// NotExpError tests that no log entry error field
// (zerolog.ErrorFieldName) is equal to exp.
func (ets Entries) NotExpError(exp string) {
	ets.t.Helper()
	ets.notExp(func(e *Entry) string { return e.expError(exp) })
}

// ExpNum tests that at least one log entry has key and its numerical
// value is equal to exp.
func (ets Entries) ExpNum(key string, exp float64) {
	ets.t.Helper()
	ets.exp(func(e *Entry) string { return e.expNum(key, exp) })
}

// NotExpNum tests that at least one log entry has key and its numerical
// value is equal to exp.
func (ets Entries) NotExpNum(key string, exp float64) {
	ets.t.Helper()
	ets.notExp(func(e *Entry) string { return e.expNum(key, exp) })
}

func (ets Entries) exp(f func(*Entry) string) {
	ets.t.Helper()
	e := ets.Get()
	for ent := range e {
		if f(e[ent]) == "" {
			return
		}
	}
	ets.t.Error("no matching log entry found")
}

// Print prints zerolog log entries.
func (ets Entries) Print() {
	ets.t.Helper()
	ets.t.Log("entries logged so far:")
	for _, e := range ets.e {
		ets.t.Log("  " + e.raw)
	}
}

func (ets Entries) notExp(f func(*Entry) string) {
	ets.t.Helper()
	e := ets.Get()
	for ent := range e {
		if f(e[ent]) == "" {
			ets.t.Error("matching log entry found")
		}
	}
}

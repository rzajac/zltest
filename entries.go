package zltest

import "time"

// Entries represents zerolog log entries.
type Entries struct {
	e []*Entry
	t T // Test manager.
}

// Get returns the list of Entry in Entries
func (ets Entries) Get() []*Entry {
	return ets.e
}

// ExpEntry returns nth logged entry.
func (ets Entries) ExpEntry(n int) *Entry {
	if n < len(ets.e) {
		return ets.e[n]
	}
	ets.t.Fatalf("expected %d%s logged entry to exist", n, ordinal(n))
	return nil
}

// ordinal returns English ordinal for a whole number.
func ordinal(n int) string {
	switch n {
	case 0:
		return "th"
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
	have := len(ets.e)
	if have != want {
		ets.t.Errorf("expected %d entries got %d", want, have)
	}
}

// ExpStr tests that at least one log entry has key, its value is a
// string, and it's equal to exp.
func (ets Entries) ExpStr(key string, exp string) {
	ets.exp(func(e *Entry) string { return e.expStr(key, exp) })
}

// NotExpStr tests that no log entry has key, its value is a
// string, and it's equal to exp.
func (ets Entries) NotExpStr(key string, exp string) {
	ets.notExp(func(e *Entry) string { return e.expStr(key, exp) })
}

// ExpTime tests that at least one log entry has key, its value is a
// string representing time in zerolog.TimeFieldFormat and it's equal
// to exp.
func (ets Entries) ExpTime(key string, exp time.Time) {
	ets.exp(func(e *Entry) string { return e.expTime(key, exp) })
}

// NotExpTime tests that no one log entry has key, its value is a
// string representing time in zerolog.TimeFieldFormat and it's equal
// to exp.
func (ets Entries) NotExpTime(key string, exp time.Time) {
	ets.notExp(func(e *Entry) string { return e.expTime(key, exp) })
}

// ExpDur tests that at least one log entry has key and its value is
// equal to exp time.Duration.  The duration vale in the entry is
// multiplied by zerolog.DurationFieldUnit before the comparison.
func (ets Entries) ExpDur(key string, exp time.Duration) {
	ets.exp(func(e *Entry) string { return e.expDur(key, exp) })
}

// NotExpDur tests that no log entry has key and its value is
// equal to exp time.Duration.  The duration vale in the entry is
// multiplied by zerolog.DurationFieldUnit before the comparison.
func (ets Entries) NotExpDur(key string, exp time.Duration) {
	ets.notExp(func(e *Entry) string { return e.expDur(key, exp) })
}

// ExpBool tests that at lest one entry has a key, its value is
// boolean and equal to exp.
func (ets Entries) ExpBool(key string, exp bool) {
	ets.exp(func(e *Entry) string { return e.expBool(key, exp) })
}

// NotExpBool tests that no log entry has a key, its value is
// boolean and equal to exp.
func (ets Entries) NotExpBool(key string, exp bool) {
	ets.notExp(func(e *Entry) string { return e.expBool(key, exp) })
}

// ExpMsg tests that at least one log entry message field
// (zerolog.MessageFieldName) is equal to exp.
func (ets Entries) ExpMsg(exp string) {
	ets.exp(func(e *Entry) string { return e.expMsg(exp) })
}

// NotExpMsg tests that no log entry message field
// (zerolog.MessageFieldName) is equal to exp.
func (ets Entries) NotExpMsg(exp string) {
	ets.notExp(func(e *Entry) string { return e.expMsg(exp) })
}

// ExpNum tests that at least one log entry has key and its numerical
// value is equal to exp.
func (ets Entries) ExpNum(key string, exp float64) {
	ets.exp(func(e *Entry) string { return e.expNum(key, exp) })
}

// NotExpNum tests that at least one log entry has key and its numerical
// value is equal to exp.
func (ets Entries) NotExpNum(key string, exp float64) {
	ets.notExp(func(e *Entry) string { return e.expNum(key, exp) })
}

func (ets Entries) exp(f func(*Entry) string) {
	e := ets.Get()
	for ent := range e {
		if f(e[ent]) == "" {
			return
		}
	}
	ets.t.Error("no matching log entry found")
}

func (ets Entries) notExp(f func(*Entry) string) {
	e := ets.Get()
	for ent := range e {
		if f(e[ent]) == "" {
			ets.t.Error("matching log entry found")
		}
	}
}

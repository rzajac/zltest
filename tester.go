package zlogtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
)

// Tester represents zerolog log entries tester.
type Tester struct {
	buf []byte     // Buffer zerolog writes to.
	tb  testing.TB // Test or benchmark manager.
}

// New creates new instance of zerolog tester.
func New(tb testing.TB) *Tester {
	return &Tester{
		buf: make([]byte, 0),
		tb:  tb,
	}
}

// Write implements io.Writer interface.
func (tst *Tester) Write(p []byte) (n int, err error) {
	tst.buf = append(tst.buf, p...)
	return len(p), nil
}

// String implements fmt.Stringer interface.
func (tst *Tester) String() string { return string(tst.buf) }

// Entries returns all logged entries. It panics if log entry cannot be decoded.
func (tst *Tester) Entries() []*Entry {
	dec := json.NewDecoder(bytes.NewReader(tst.buf))
	ets := make([]*Entry, 0)
	for {
		m := make(map[string]interface{})
		if err := dec.Decode(&m); err != nil {
			if errors.Is(err, io.EOF) {
				return ets
			}
			panic(err)
		}
		ets = append(ets, &Entry{m, tst.tb})
	}
	return ets
}

// FirstEntry returns first log entry.
func (tst *Tester) FirstEntry() *Entry {
	ets := tst.Entries()
	if len(ets) == 0 {
		return nil
	}
	return ets[0]
}

// LastEntry returns last log entry.
func (tst *Tester) LastEntry() *Entry {
	ets := tst.Entries()
	if len(ets) == 0 {
		return nil
	}
	return ets[len(ets)-1]
}

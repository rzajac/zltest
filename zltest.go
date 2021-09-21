// Package zltest provides facilities to test zerolog log messages.
package zltest

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/rs/zerolog"
)

// Tester represents zerolog log tester.
type Tester struct {
	mx  sync.RWMutex // Guards the buffer.
	buf []byte       // Buffer zerolog writes to.
	cnt int          // Number of all log messages (calls to Write).
	t   T            // Test manager.
}

// New creates new instance of zerolog tester.
func New(t T) *Tester {
	return &Tester{
		buf: make([]byte, 0, 500),
		t:   t,
	}
}

// Logger returns zerolog.Logger using this tester as io.Writer.
func (tst *Tester) Logger() zerolog.Logger {
	return zerolog.New(tst)
}

// Write implements io.Writer interface.
func (tst *Tester) Write(p []byte) (n int, err error) {
	tst.mx.Lock()
	defer tst.mx.Unlock()

	tst.cnt++
	tst.buf = append(tst.buf, p...)
	return len(p), nil
}

// Len returns number of log messages written to the Tester.
func (tst *Tester) Len() int {
	return tst.cnt
}

// String implements fmt.Stringer interface and returns everything written
// to the Tester so far. Calls Fatal on error.
func (tst *Tester) String() string {
	tst.mx.RLock()
	defer tst.mx.RUnlock()
	return string(tst.buf)
}

// Entries returns all logged entries in the order they were logged. It calls
// Fatal if any of the log entries cannot be decoded or Scanner error.
func (tst *Tester) Entries() Entries {
	tst.mx.RLock()
	defer tst.mx.RUnlock()
	tst.t.Helper()

	ets := make([]*Entry, 0, tst.cnt)

	var off int64
	dec := json.NewDecoder(bytes.NewReader(tst.buf))
	for dec.More() {
		m := make(map[string]interface{})
		if err := dec.Decode(&m); err != nil {
			tst.t.Fatal(err)
			return Entries{t: tst.t}
		}

		tmp := tst.buf[off:dec.InputOffset()]
		off = dec.InputOffset()
		ets = append(ets, &Entry{
			raw: string(bytes.TrimSpace(tmp)),
			m:   m,
			t:   tst.t,
		})
	}

	return Entries{e: ets, t: tst.t}
}

// Filter returns only entries matching log level.
func (tst *Tester) Filter(level zerolog.Level) Entries {
	ets := make([]*Entry, 0)
	for _, ent := range tst.Entries().Get() {
		if lvl, _ := ent.Str(zerolog.LevelFieldName); lvl == level.String() {
			ets = append(ets, ent)
		}
	}
	return Entries{e: ets, t: tst.t}
}

// FirstEntry returns first log entry or nil if no log entries written
// to the Tester. It calls Fatal if any of the log entries cannot be decoded.
func (tst *Tester) FirstEntry() *Entry {
	tst.mx.RLock()
	defer tst.mx.RUnlock()

	ets := tst.Entries().Get()
	if len(ets) == 0 {
		return nil
	}
	return ets[0]
}

// LastEntry returns last log entry or nil if no log entries written
// to the Tester. It calls Fatal if any of the log entries cannot be decoded.
func (tst *Tester) LastEntry() *Entry {
	tst.mx.RLock()
	defer tst.mx.RUnlock()

	ets := tst.Entries().Get()
	if len(ets) == 0 {
		return nil
	}
	return ets[len(ets)-1]
}

// Reset resets the Tester.
func (tst *Tester) Reset() {
	tst.mx.Lock()
	defer tst.mx.Unlock()

	tst.cnt = 0
	tst.buf = tst.buf[:0]
}

// T is a subset of testing.TB interface.
// It's primarily used to test zltest package but can be used to implement
// custom actions to be taken on errors.
type T interface {
	// Error is equivalent to Log followed by Fail.
	Error(args ...interface{})

	// Errorf is equivalent to Logf followed by Fail.
	Errorf(format string, args ...interface{})

	// Fatal is equivalent to Log followed by FailNow.
	Fatal(args ...interface{})

	// Fatalf is equivalent to Logf followed by FailNow.
	Fatalf(format string, args ...interface{})

	// Helper marks the calling function as a test helper function.
	// When printing file and line information, that function will be skipped.
	// Helper may be called simultaneously from multiple goroutines.
	Helper()

	// Log formats its arguments using default formatting, analogous to Println,
	// and records the text in the error log. For tests, the text will be printed only if
	// the test fails or the -test.v flag is set. For benchmarks, the text is always
	// printed to avoid having performance depend on the value of the -test.v flag.
	Log(args ...interface{})

	// Logf formats its arguments according to the format, analogous to Printf, and
	// records the text in the error log. A final newline is added if not provided. For
	// tests, the text will be printed only if the test fails or the -test.v flag is
	// set. For benchmarks, the text is always printed to avoid having performance
	// depend on the value of the -test.v flag.
	Logf(format string, args ...interface{})

	// Cleanup registers a function to be called when the test (or subtest) and all its
	// subtests complete. Cleanup functions will be called in last added,
	// first called order.
	Cleanup(func())
}

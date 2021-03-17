package zltest

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	. "github.com/rzajac/zltest/internal"
)

func Test_Entry_String(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()

	// --- When ---
	got := tst.LastEntry().String()

	// --- Then ---
	exp := `{"level":"error","key0":"val0"}`
	assert.Exactly(t, exp, got)
}

func Test_Entry_Str(t *testing.T) {
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Str("str", "val").Int("int", 42).Send()

	tt := []struct {
		testN string

		key    string
		expVal string
		expSt  KeyStatus
	}{
		{"1", "str", "val", KeyFound},
		{"2", "int", "", KeyBadType},
		{"3", "missing", "", KeyMissing},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			val, st := tst.LastEntry().Str(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.expVal, val, "test %s", tc.testN)
			assert.Exactly(t, tc.expSt, st, "test %s", tc.testN)
		})
	}
}

func Test_Entry_Float64(t *testing.T) {
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Float64("float64", 12.3).Str("str", "val").Send()

	tt := []struct {
		testN string

		key    string
		expVal float64
		expSt  KeyStatus
	}{
		{"1", "float64", 12.3, KeyFound},
		{"2", "str", 0.0, KeyBadType},
		{"3", "missing", 0.0, KeyMissing},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			val, st := tst.LastEntry().Float64(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.expVal, val, "test %s", tc.testN)
			assert.Exactly(t, tc.expSt, st, "test %s", tc.testN)
		})
	}
}

func Test_Entry_Bool(t *testing.T) {
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Bool("bool", true).Str("str", "val").Send()

	tt := []struct {
		testN string

		key    string
		expVal bool
		expSt  KeyStatus
	}{
		{"1", "bool", true, KeyFound},
		{"2", "str", false, KeyBadType},
		{"3", "missing", false, KeyMissing},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			val, st := tst.LastEntry().Bool(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.expVal, val, "test %s", tc.testN)
			assert.Exactly(t, tc.expSt, st, "test %s", tc.testN)
		})
	}
}

func Test_Entry_Time(t *testing.T) {
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	now := time.Date(2020, 11, 18, 22, 17, 4, 948442004, time.UTC)
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Time("time", now).Str("str", "val").Int("int", 42).Send()

	tt := []struct {
		testN string

		key    string
		expVal time.Time
		expSt  KeyStatus
	}{
		{"1", "time", now, KeyFound},
		{"2", "str", time.Time{}, KeyBadFormat},
		{"3", "int", time.Time{}, KeyBadType},
		{"4", "missing", time.Time{}, KeyMissing},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			val, st := tst.LastEntry().Time(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.expVal, val, "test %s", tc.testN)
			assert.Exactly(t, tc.expSt, st, "test %s", tc.testN)
		})
	}
}

func Test_Entry(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Error().Int("key0", 123).Msg("message")

	// --- Then ---
	entry := tst.LastEntry()

	entry.ExpNum("key0", 123)
	entry.ExpMsg("message")
	entry.ExpLevel(zerolog.ErrorLevel)
	fmt.Println(tst.String())
}

func Test_Entry_ExpStr_Equal(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpStr("key", "val")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpStr_NotEqual(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry key 'key' to have value 'value' but got 'val'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpStr("key", "value")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpStr_NotFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry to have key 'some_key'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpStr("some_key", "value")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpTime_Equal(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")

	now := time.Now()
	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Time("key", now).Send()

	// --- When ---
	tst.LastEntry().ExpTime("key", now)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpTime_NotEqual(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	exp := time.Now()
	got := exp.Add(time.Second)

	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		fmt.Sprintf("expected entry '%s' to be '%s' but is '%s'",
			"key",
			exp.Format(time.RFC3339Nano),
			got.Format(time.RFC3339Nano),
		),
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Time("key", got).Send()

	// --- When ---
	tst.LastEntry().ExpTime("key", exp)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpTime_NotFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry to have key 'some_key'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpTime("some_key", time.Time{})

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpDur_Equal(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	dur := 42 * time.Second
	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Dur("key", dur).Send()

	// --- When ---
	tst.LastEntry().ExpDur("key", dur)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpDur_NotEqual(t *testing.T) {
	// --- Given ---
	exp := 42 * time.Second
	got := 44 * time.Second

	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		fmt.Sprintf("expected entry key '%s' to have value '%d' (%s) but got '%d' (%s)",
			"key",
			exp/zerolog.DurationFieldUnit,
			exp.String(),
			got/zerolog.DurationFieldUnit,
			got.String(),
		),
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Dur("key", got).Send()

	// --- When ---
	tst.LastEntry().ExpDur("key", exp)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpDur_NotFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry to have key 'some_key'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpDur("some_key", time.Second)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpBool_Equal(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Bool("key", true).Send()

	// --- When ---
	tst.LastEntry().ExpBool("key", true)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpBool_NotEqual(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry key 'key' to have value 'true' but got 'false'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Bool("key", false).Send()

	// --- When ---
	tst.LastEntry().ExpBool("key", true)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpBool_NotFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry to have key 'some_key'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Bool("key", true).Send()

	// --- When ---
	tst.LastEntry().ExpBool("some_key", true)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpTimeWithin_Equal(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")

	got := time.Now()
	exp := got.Add(time.Second)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Time("key", got).Send()

	// --- When ---
	tst.LastEntry().ExpTimeWithin("key", exp, 2*time.Second)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpTimeWithin_NotEqual(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	exp := time.Now()
	got := exp.Add(3 * time.Second)

	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Errorf",
		"expected entry '%s' to be within '%s' but difference is '%s'",
		"key",
		"2s",
		"3s",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Time("key", got).Send()

	// --- When ---
	tst.LastEntry().ExpTimeWithin("key", exp, 2*time.Second)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpTimeWithin_NotFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry to have key 'some_key'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpTimeWithin("some_key", time.Time{}, time.Second)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpNum_Equal(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Float64("key", 1.23).Send()

	// --- When ---
	tst.LastEntry().ExpNum("key", 1.23)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpNum_NotEqual(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry key 'key' to have value '1.231' but got '1.23'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Float64("key", 1.23).Send()

	// --- When ---
	tst.LastEntry().ExpNum("key", 1.231)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpNum_NotFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry to have key 'some_key'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Float64("key", 1.23).Send()

	// --- When ---
	tst.LastEntry().ExpNum("some_key", 1.23)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpLoggedWithin_Equal(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst).With().Timestamp().Logger()
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpLoggedWithin(time.Now(), time.Millisecond)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpLoggedWithin_NotEqual(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Errorf",
		"expected entry '%s' to be within '%s' but difference is '%s'",
		"time",
		"1µs",
		mock.Anything,
	)

	tst := New(mck)
	log := zerolog.New(tst).With().Timestamp().Logger()
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpLoggedWithin(time.Now(), time.Microsecond)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entry_ExpLoggedWithin_NotFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On(
		"Error",
		"expected entry to have key 'some_key'",
	)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key", "val").Send()

	// --- When ---
	tst.LastEntry().ExpTimeWithin("some_key", time.Time{}, time.Second)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_formatError(t *testing.T) {
	tt := []struct {
		testN string

		status KeyStatus
		key    string
		typ    string
		exp    string
	}{
		{"1", KeyMissing, "key", "number", "expected entry to have key 'key'"},
		{"2", KeyBadType, "key", "number", "expected entry key 'key' to be 'number'"},
		{"3", KeyBadFormat, "key", "number", "key 'key' in a wrong format"},
		{"4", "unknown", "key", "number", "invalid KeyStatus 'unknown'"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			got := formatError(tc.status, tc.key, tc.typ)

			// --- Then ---
			assert.Exactly(t, tc.exp, got, "test %s", tc.testN)
		})
	}
}

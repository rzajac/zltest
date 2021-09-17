package zltest

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	. "github.com/rzajac/zltest/internal"
)

func Test_Entries(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Error().Int("key0", 0).Msg("message0")
	log.Error().Int("key1.1", 1).Bool("key1.2", true).Msg("message1")
	log.Error().Float64("key2", 2.2).Msg("message2")

	// --- Then ---
	tst.Entries().ExpNum("key2", 2.2)
	tst.Entries().ExpNum("key0", 0)
	tst.Entries().ExpMsg("message1")
	tst.Entries().ExpNum("key1.1", 1)
	tst.Entries().ExpBool("key1.2", true)
	tst.Entries().ExpMsg("message0")
}

func Test_Entries_ExpEntry(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Int("key0", 0).Msg("message0")
	log.Error().Int("key1.1", 1).Bool("key1.2", true).Msg("message1")
	log.Error().Float64("key2", 2.2).Msg("message2")

	// --- Then ---
	ent := tst.Entries().ExpEntry(0)
	ent.ExpNum("key0", 0)

	ent = tst.Entries().ExpEntry(1)
	ent.ExpNum("key1.1", 1)
	ent.ExpBool("key1.2", true)
	ent.ExpMsg("message1")

	ent = tst.Entries().ExpEntry(2)
	ent.ExpNum("key2", 2.2)
	ent.ExpMsg("message2")
}

func Test_Entries_ExpEntry_fatal(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Fatalf", "expected %d%s logged entry to exist", 3, "rd")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Int("key0", 0).Msg("message0")
	log.Error().Int("key1.1", 1).Bool("key1.2", true).Msg("message1")
	log.Error().Float64("key2", 2.2).Msg("message2")

	// --- When ---
	ent := tst.Entries().ExpEntry(3)

	// --- Then ---
	mck.AssertExpectations(t)
	assert.Nil(t, ent)
}

func Test_ordinal(t *testing.T) {
	tt := []struct {
		testN string

		num int
		exp string
	}{
		{"1", -1, "th"},
		{"2", 0, "th"},
		{"3", 1, "st"},
		{"4", 2, "nd"},
		{"5", 3, "rd"},
		{"6", 4, "th"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			assert.Exactly(t, tc.exp, ordinal(tc.num))
		})
	}
}

func Test_Entries_ExpLen(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Error().Int("key0", 0).Msg("message0")
	log.Error().Int("key1.1", 1).Bool("key1.2", true).Msg("message1")
	log.Error().Float64("key2", 2.2).Msg("message2")

	// --- Then ---
	tst.Entries().ExpLen(3)
}

func Test_Entries_ExpLen_error(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Errorf", "expected %d entries got %d", 3, 1)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Int("key0", 0).Msg("message0")

	// --- When ---
	tst.Entries().ExpLen(3)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_NotExpMsg_notFound(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Error().Int("key0", 0).Msg("message0")
	log.Error().Int("key1.1", 1).Bool("key1.2", true).Msg("message1")
	log.Error().Float64("key2", 2.2).Msg("message2")

	// --- Then ---
	tst.Entries().NotExpMsg("message")
}

func Test_Entries_ExpBool_empty(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	// --- When ---
	New(mck).Entries().ExpBool("key1", false)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpBool_found(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Bool("key", true).Send()
	log.Error().Bool("key", false).Send()

	// --- When ---
	tst.Entries().ExpBool("key", true)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpBool_notFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Bool("key", true).Send()
	log.Error().Bool("key", true).Send()

	// --- When ---
	tst.Entries().ExpBool("key", false)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpTime_found(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")

	now := time.Now()
	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Time("key", time.Now()).Send()
	log.Error().Time("key", now).Send()
	log.Error().Time("key", time.Now()).Send()

	// --- When ---
	tst.Entries().ExpTime("key", now)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpTime_notFound(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	exp := time.Now()
	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Time("key", time.Now()).Send()
	log.Error().Time("key", exp.Add(time.Second)).Send()
	log.Error().Time("key", time.Now()).Send()

	// --- When ---
	tst.Entries().ExpTime("key", exp)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpDur_found(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")

	dur := 42 * time.Second

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Dur("key", 43*time.Second).Send()
	log.Error().Time("key", time.Now()).Send()
	log.Error().Dur("key", dur).Send()

	// --- When ---
	tst.Entries().ExpDur("key", dur)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpDur_notFound(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Dur("key", 43*time.Second).Send()
	log.Error().Time("key", time.Now()).Send()
	log.Error().Dur("key", 43*time.Second).Send()

	// --- When ---
	tst.Entries().ExpDur("key", 42*time.Second)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpNum_found(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Float64("key", 1.23).Send()
	log.Error().Float64("key", 0).Send()
	log.Error().Int("key", 42).Send()
	log.Error().Float64("key", -1).Send()

	// --- When ---
	tst.Entries().ExpNum("key", 1.23)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpNum_notFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Float64("key", 1.22).Send()
	log.Error().Float64("key", 0).Send()
	log.Error().Int("key", 42).Send()
	log.Error().Float64("key", -1).Send()

	// --- When ---
	tst.Entries().ExpNum("key", 1.23)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpStr_found(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Error().Str("key2", "val2").Send()

	// --- When ---
	tst.Entries().ExpStr("key1", "val1")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpStr_foundFirst(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Error().Str("key2", "val2").Send()

	// --- When ---
	tst.Entries().ExpStr("key0", "val0")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpStr_foundLast(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Error().Str("key2", "val2").Send()

	// --- When ---
	tst.Entries().ExpStr("key2", "val2")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpStr_notFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Error().Str("key2", "val2").Send()

	// --- When ---
	tst.Entries().ExpStr("key1", "val")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpStr_filterFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Debug().Str("key2", "val2").Send()

	// --- When ---
	tst.Filter(zerolog.InfoLevel).ExpStr("key1", "val1")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpStr_filterNotFound(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Debug().Str("key2", "val2").Send()

	// --- When ---
	tst.Filter(zerolog.InfoLevel).ExpStr("key2", "val2")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpStr_noKey(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Error().Str("key2", "val2").Send()

	// --- When ---
	tst.Entries().ExpStr("key", "val0")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_NotExpStr(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Error().Str("key2", "val2").Send()

	// --- When ---
	tst.Entries().NotExpStr("key", "val0")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_NotExpStr_found(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "matching log entry found")

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Error().Str("key2", "val2").Send()

	// --- When ---
	tst.Entries().NotExpStr("key1", "val1")

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_Print(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Log", "entries logged so far:")
	mck.On("Log", `  {"level":"error","key0":"val0"}`)
	mck.On("Log", `  {"level":"error","key1":"val1"}`)

	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()

	// --- When ---
	tst.Entries().Print()

	// --- Then ---
	mck.AssertExpectations(t)
}

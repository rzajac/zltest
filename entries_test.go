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

func Test_Entries_ExpEntry_Fatal(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
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

func Test_Entries_Len(t *testing.T) {
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

func Test_Entries_Len_Error(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Errorf", "expected %d entries got %d", 3, 1)

	tst := New(mck)
	log := zerolog.New(tst)

	log.Error().Int("key0", 0).Msg("message0")

	// --- When ---
	tst.Entries().ExpLen(3)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_NotFound(t *testing.T) {
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

func Test_Entries_Empty(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Error", "no matching log entry found")
	tst := New(mck)

	// --- When ---
	tst.Entries().ExpBool("key1", false)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpBool_Found(t *testing.T) {
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
func Test_Entries_ExpBool_NotFound(t *testing.T) {
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

func Test_Entries_ExpTime_Found(t *testing.T) {
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

func Test_Entries_ExpTime_NotFound(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	exp := time.Now()
	got := exp.Add(time.Second)
	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Time("key", time.Now()).Send()
	log.Error().Time("key", got).Send()
	log.Error().Time("key", time.Now()).Send()

	// --- When ---
	tst.Entries().ExpTime("key", exp)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpDur_Found(t *testing.T) {
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

func Test_Entries_ExpDur_NotFound(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	mck := &TMock{}
	mck.On("Helper")
	mck.On("Error", "no matching log entry found")

	got := 43 * time.Second
	exp := 42 * time.Second
	tst := New(mck)
	log := zerolog.New(tst)
	log.Error().Dur("key", 43*time.Second).Send()
	log.Error().Time("key", time.Now()).Send()
	log.Error().Dur("key", got).Send()

	// --- When ---
	tst.Entries().ExpDur("key", exp)

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Entries_ExpNum_Found(t *testing.T) {
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

func Test_Entries_ExpNum_NotFound(t *testing.T) {
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

func Test_Entries_ExpString_Found(t *testing.T) {
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

func Test_Entries_ExpString_Found_First(t *testing.T) {
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

func Test_Entries_ExpString_Found_Last(t *testing.T) {
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

func Test_Entries_ExpString_NotFound(t *testing.T) {
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

func Test_Entries_ExpString_Filter_Found(t *testing.T) {
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

func Test_Entries_ExpString_Filter_NotFound(t *testing.T) {
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

func Test_Entries_ExpString_NoKey(t *testing.T) {
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

func Test_Entries_ExpString_NotExp(t *testing.T) {
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

func Test_Entries_ExpString_NotExp_Found(t *testing.T) {
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

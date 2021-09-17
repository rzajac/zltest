package zltest

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	. "github.com/rzajac/zltest/internal"
)

func Test_Tester_Logger(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := tst.Logger()

	// --- When ---
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()

	// --- Then ---
	assert.Exactly(t, 2, tst.Len())
}

func Test_Tester_Len(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()

	// --- Then ---
	assert.Exactly(t, 2, tst.Len())
}

func Test_Tester_String(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()

	// --- Then ---
	exp := "{\"level\":\"info\",\"key0\":\"val0\"}\n{\"level\":\"error\",\"key1\":\"val1\"}\n"
	assert.Exactly(t, exp, tst.String())
}

func Test_Tester_Entries(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()

	// --- Then ---
	assert.Len(t, tst.Entries().Get(), 2)
}

func Test_Tester_Filter(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()
	log.Debug().Str("keyD", "valD").Send()

	// --- Then ---
	assert.Len(t, tst.Filter(zerolog.InfoLevel).Get(), 1)
	assert.Len(t, tst.Filter(zerolog.ErrorLevel).Get(), 1)
	assert.Len(t, tst.Filter(zerolog.DebugLevel).Get(), 1)
	assert.Len(t, tst.Filter(zerolog.FatalLevel).Get(), 0)
}

func Test_Tester_Entries_noEntries(t *testing.T) {
	// --- Given ---
	tst := New(t)

	// --- When ---
	ets := tst.Entries().Get()

	// --- Then ---
	assert.NotNil(t, ets)
	assert.Len(t, ets, 0)
}

func Test_Tester_Entries_errorDecodingEntry(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
	mck.On("Helper")
	mck.On("Fatal", mock.AnythingOfType("*json.SyntaxError"))

	tst := New(mck)

	// --- When ---
	_, _ = tst.Write([]byte("{ bad json }"))
	tst.Entries()

	// --- Then ---
	mck.AssertExpectations(t)
}

func Test_Tester_FirstEntry(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()

	// --- When ---
	got := tst.FirstEntry().String()

	// --- Then ---
	exp := `{"level":"info","key0":"val0"}`
	assert.Exactly(t, exp, got)
}

func Test_Tester_FirstEntry_noEntries(t *testing.T) {
	// --- Given ---
	tst := New(t)

	// --- Then ---
	assert.Nil(t, tst.FirstEntry())
}

func Test_Tester_LastEntry(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()

	// --- When ---
	got := tst.LastEntry().String()

	// --- Then ---
	exp := `{"level":"error","key1":"val1"}`
	assert.Exactly(t, exp, got)
}

func Test_Tester_LastEntry_noEntries(t *testing.T) {
	// --- Given ---
	tst := New(t)

	// --- Then ---
	assert.Nil(t, tst.LastEntry())
}

func Test_Tester_Reset(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Info().Str("key0", "val0").Send()
	log.Error().Str("key1", "val1").Send()

	// --- When ---
	tst.Reset()

	// --- Then ---
	assert.Exactly(t, 0, tst.Len())
	assert.Exactly(t, "", tst.String())
}

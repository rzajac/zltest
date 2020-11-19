package zltest

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	. "github.com/rzajac/zltest/internal"
)

func Test_Tester_Len(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Info().Str("keyI", "valI").Send()
	log.Error().Str("keyE", "valE").Send()

	// --- Then ---
	assert.Exactly(t, 2, tst.Len())
}

func Test_Tester_String(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Info().Str("keyI", "valI").Send()
	log.Error().Str("keyE", "valE").Send()

	// --- Then ---
	exp := "{\"level\":\"info\",\"keyI\":\"valI\"}\n{\"level\":\"error\",\"keyE\":\"valE\"}\n"
	assert.Exactly(t, exp, tst.String())
}

func Test_Tester_Entries(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Info().Str("keyI", "valI").Send()
	log.Error().Str("keyE", "valE").Send()

	// --- Then ---
	assert.Len(t, tst.Entries(), 2)
}

func Test_Tester_Filter(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)

	// --- When ---
	log.Info().Str("keyI", "valI").Send()
	log.Error().Str("keyE", "valE").Send()
	log.Debug().Str("keyD", "valD").Send()

	// --- Then ---
	assert.Len(t, tst.Filter(zerolog.InfoLevel), 1)
	assert.Len(t, tst.Filter(zerolog.ErrorLevel), 1)
	assert.Len(t, tst.Filter(zerolog.DebugLevel), 1)
	assert.Len(t, tst.Filter(zerolog.FatalLevel), 0)
}

func Test_Tester_Entries_NoEntries(t *testing.T) {
	// --- Given ---
	tst := New(t)

	// --- When ---
	ets := tst.Entries()

	// --- Then ---
	assert.NotNil(t, ets)
	assert.Len(t, ets, 0)
}

func Test_Tester_Entries_ErrorDecodingEntry(t *testing.T) {
	// --- Given ---
	mck := &TMock{}
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
	log.Info().Str("keyI", "valI").Send()
	log.Error().Str("keyE", "valE").Send()

	// --- When ---
	got := tst.FirstEntry().String()

	// --- Then ---
	exp := `{"level":"info","keyI":"valI"}`
	assert.Exactly(t, exp, got)
}

func Test_Tester_FirstEntry_NoEntries(t *testing.T) {
	// --- Given ---
	tst := New(t)

	// --- Then ---
	assert.Nil(t, tst.FirstEntry())
}

func Test_Tester_LastEntry(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Info().Str("keyI", "valI").Send()
	log.Error().Str("keyE", "valE").Send()

	// --- When ---
	got := tst.LastEntry().String()

	// --- Then ---
	exp := `{"level":"error","keyE":"valE"}`
	assert.Exactly(t, exp, got)
}

func Test_Tester_LastEntry_NoEntries(t *testing.T) {
	// --- Given ---
	tst := New(t)

	// --- Then ---
	assert.Nil(t, tst.LastEntry())
}

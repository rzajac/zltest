package zlogtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_Entry_Str_KeyFound(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()

	// --- When ---
	got, status := tst.LastEntry().Str("key0")

	// --- Then ---
	assert.Exactly(t, KeyFound, status)
	assert.Exactly(t, "val0", got)
}

func Test_Entry_Str_KeyBadType(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Int("key0", 0).Send()

	// --- When ---
	got, status := tst.LastEntry().Str("key0")

	// --- Then ---
	assert.Exactly(t, KeyBadType, status)
	assert.Exactly(t, "", got)
}

func Test_Entry_Str_KeyMissing(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Send()

	// --- When ---
	got, status := tst.LastEntry().Str("key0")

	// --- Then ---
	assert.Exactly(t, KeyMissing, status)
	assert.Exactly(t, "", got)
}

func Test_Entry_Float64_KeyFound(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Float64("key0", 1.1).Send()

	// --- When ---
	got, status := tst.LastEntry().Float64("key0")

	// --- Then ---
	assert.Exactly(t, KeyFound, status)
	assert.Exactly(t, 1.1, got)
}

func Test_Entry_Float64_KeyBadType(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()

	// --- When ---
	got, status := tst.LastEntry().Float64("key0")

	// --- Then ---
	assert.Exactly(t, KeyBadType, status)
	assert.Exactly(t, 0.0, got)
}

func Test_Entry_Float64_KeyMissing(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Send()

	// --- When ---
	got, status := tst.LastEntry().Float64("key0")

	// --- Then ---
	assert.Exactly(t, KeyMissing, status)
	assert.Exactly(t, 0.0, got)
}

func Test_Entry_Bool_KeyFound(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Bool("key0", true).Send()

	// --- When ---
	got, status := tst.LastEntry().Bool("key0")

	// --- Then ---
	assert.Exactly(t, KeyFound, status)
	assert.Exactly(t, true, got)
}

func Test_Entry_Bool_KeyBadType(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()

	// --- When ---
	got, status := tst.LastEntry().Bool("key0")

	// --- Then ---
	assert.Exactly(t, KeyBadType, status)
	assert.Exactly(t, false, got)
}

func Test_Entry_Bool_KeyMissing(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Send()

	// --- When ---
	got, status := tst.LastEntry().Bool("key0")

	// --- Then ---
	assert.Exactly(t, KeyMissing, status)
	assert.Exactly(t, false, got)
}

func Test_Entry_Time_KeyFound(t *testing.T) {
	// --- Given ---
	old := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = time.RFC3339Nano
	defer func() { zerolog.TimeFieldFormat = old }()

	now := time.Now()
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Time("key0", now).Send()

	// --- When ---
	got, status := tst.LastEntry().Time("key0")

	// --- Then ---
	assert.Exactly(t, KeyFound, status)
	assert.Exactly(t, now.Format(time.RFC3339Nano), got.Format(time.RFC3339Nano))
}

func Test_Entry_Time_KeyBadFormat(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Str("key0", "val0").Send()

	// --- When ---
	got, status := tst.LastEntry().Time("key0")

	// --- Then ---
	assert.Exactly(t, KeyBadFormat, status)
	assert.Exactly(t, time.Time{}, got)
}

func Test_Entry_Time_KeyBadType(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Int("key0", 123).Send()

	// --- When ---
	got, status := tst.LastEntry().Time("key0")

	// --- Then ---
	assert.Exactly(t, KeyBadType, status)
	assert.Exactly(t, time.Time{}, got)
}

func Test_Entry_Time_KeyMissing(t *testing.T) {
	// --- Given ---
	tst := New(t)
	log := zerolog.New(tst)
	log.Error().Send()

	// --- When ---
	got, status := tst.LastEntry().Time("key0")

	// --- Then ---
	assert.Exactly(t, KeyMissing, status)
	assert.Exactly(t, time.Time{}, got)
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

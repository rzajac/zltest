# zlogtest

Package `zlogtest` provides facilities to test [zerolog](https://github.com/rs/zerolog) 
log messages.

```
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
```

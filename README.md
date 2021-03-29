[![Go Report Card](https://goreportcard.com/badge/github.com/rzajac/zltest)](https://goreportcard.com/report/github.com/rzajac/zltest)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rzajac/zltest)

Logging is an integral part of most applications and as such it has tobe tested. 
Package `zltest` provides facilities to test 
[zerolog](https://github.com/rs/zerolog) log messages.

In general `zltest.Tester` provides methods to test values (or existence) of
specific fields in logged messages.

## Installation

```
go get github.com/rzajac/zltest
```

## Examples

```go
func Test_ServiceLogsProperly(t *testing.T) {
    // --- Given ---
    // Crate zerolog test helper. 
    // tst := zltest.New(t)

    // Configure zerolog and pas tester as a writer.     
    log := zerolog.New(tst).With().Timestamp().Logger()
    
    // Inject log to tested service or package.
    srv := MyService(log)

    // --- When ---
    srv.ExecuteSomeLogic()

    // --- Then ---

    // Test if log messages were generated properly.
    ent := tst.LastEntry()
    ent.ExpNum("key0", 123)
    ent.ExpMsg("message")
    ent.ExpLevel(zerolog.ErrorLevel)
}
```

## License

BSD-2-Clause
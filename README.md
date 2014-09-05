sdk-go
======

SDK in Go for THECALLR API

## Quick start
Install the Go way

    go get github.com/THECALLR/sdk-go

Or get sources from Github

## Initialize your code

```go
import "github.com/THECALLR/sdk-go"
```

## Basic Example (Send SMS)
See full example in [samples/quickstart.go](samples/quickstart.go)

```go
// Set your cedentials
// third param nil is an optional configuration. see samples
thecallr.Setup("login", "password", nil)

// 1. "call" method: each parameter of the method as an argument
result, err := thecallr.Call("sms.send", "THECALLR", "+33123456789", "Hello, world", map[string]interface{}{
	"flash_message": false,
})

// 2. "send" method: parameter of the method is an array
my_array := []interface{}{
	"THECALLR",
	"+33123456789",
	"Hello, world",
	map[string]interface{}{
		"flash_message": false,
	},
}
result, err = thecallr.Send("sms.send", my_array)
```

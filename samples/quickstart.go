// run this file like this:
// $ CALLR_API_KEY=yourapikey go run quickstart.go +15559820800
// replace the number with your personal number and you should receive the SMS.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	callr "github.com/THECALLR/sdk-go/v2"
)

func main() {
	// optional: use slog instead of log.Printf
	callr.SetLogFunc(func(format string, args ...any) {
		slog.Warn(fmt.Sprintf(strings.TrimPrefix(format, "[warning] "), args...))
	})

	// Api Key Auth (use the customer portal to generate keys)
	api := callr.NewWithAPIKeyAuth(os.Getenv("CALLR_API_KEY"))

	// optional: set a proxy
	// proxy must be in url standard format
	// http[s]://user:password@host:port
	// http[s]://host:port
	// http[s]://host
	// api.SetProxy("http://proxy:port")

	// check for destination phone number parameter
	if len(os.Args) < 2 {
		// fmt.Println("Please supply destination phone number!")
		slog.Error("Please supply destination phone number!")
		os.Exit(1)
	}

	// our context
	ctx := context.Background()

	// Send a SMS with "sms.send" JSON-RPC method
	result, err := api.Call(ctx, "sms.send", "SMS", os.Args[1], "Hello, world", nil)

	// error management
	if err != nil {
		switch e := err.(type) {
		case *callr.JSONRPCError:
			slog.Error("JSON-RPC Error",
				"code", e.Code,
				"message", e.Message,
				"data", e.Data)
		case *callr.HTTPError:
			slog.Error("HTTP Error",
				"code", e.Code,
				"message", e.Message)
		default:
			slog.Error("Other error", "error", err)
		}
		os.Exit(1)
	}

	// the sms.send JSON-RPC method returns a string
	var hash string

	if err := json.Unmarshal(result, &hash); err != nil {
		slog.Error("Error unmarshalling result", "error", err)
		os.Exit(1)
	}

	slog.Info("SMS sent", "hash", hash)
}

// run this file like this:
// $ CALLR_API_KEY=yourapikey go run quickstart.go +15559820800
// obviously, replace the number with your personal number
// and you should receive an SMS.

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	callr "github.com/THECALLR/sdk-go/v2"
)

func main() {
	// use Api Key Auth (recommended) - use the customer portal to generate keys
	api := callr.NewWithAPIKeyAuth(os.Getenv("CALLR_API_KEY"))

	// optional: set a proxy
	// proxy must be in url standard format
	// http[s]://user:password@host:port
	// http[s]://host:port
	// http[s]://host
	// api.SetProxy("http://proxy:port")

	// check for destination phone number parameter
	if len(os.Args) < 2 {
		fmt.Println("Please supply destination phone number!")
		os.Exit(1)
	}

	// context
	ctx := context.Background()

	// Example to send a SMS
	// 1. "call" method: each parameter of the method as an argument
	result, err := api.Call(ctx, "sms.send", "SMS", os.Args[1], "Hello, world", nil)

	// error management
	if err != nil {
		var jsonRpcError *callr.JSONRPCError
		if errors.As(err, &jsonRpcError) {
			fmt.Printf("Remote error: code:%d message:%s data:%v\n",
				jsonRpcError.Code, jsonRpcError.Message, jsonRpcError.Data)
		} else {
			fmt.Println("Transport error: ", err)
		}
		os.Exit(1)
		return
	}

	fmt.Println(result)
}

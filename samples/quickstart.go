package main

import (
	"fmt"
	"os"
	"github.com/THECALLR/sdk-go"
)

func main() {

	// initialize instance Callr
	// retrieve CALLR credentials from environment variables
	callr.Auth = LoginPasswordAuth(os.Getenv("CALLR_LOGIN"), os.Getenv("CALLR_PASS"))

	// if you need to set a proxy, you have to call the deprecated Setup method or directly
	// feed the callr.Config.Proxy attribute
	// http[s]://user:password@host:port
	// http[s]://host:port
	// http[s]://host

	// var config callr.Config
	// config.Proxy = "http://foo:bar@example.com:8080"

	// check for destination phone number parameter
	if len(os.Args) < 2 {
		fmt.Println("Please supply destination phone number!")
		os.Exit(1)
	}

	// Example to send a SMS
	// 1. "call" method: each parameter of the method as an argument
	result, err := callr.Call("sms.send", "SMS", os.Args[1], "Hello, world", map[string]interface{}{
		"flash_message": false,
	})

	// error management
	if err != nil {
		fmt.Println("Code:", err.Code)
		fmt.Println("Message:", err.Msg)
		fmt.Println("Data:", err.Data)
		os.Exit(1)
	} else {
		fmt.Println(result)
	}

	// 2. "send" method: parameter of the method is an array
	my_array := []interface{}{
		"SMS",
		 os.Args[1],
		"Hello, world",
		map[string]interface{}{
			"flash_message": false,
		},
	}
	result, err = callr.Send("sms.send", my_array)

	// error management
	if err != nil {
		fmt.Println("Code:", err.Code)
		fmt.Println("Message:", err.Msg)
		fmt.Println("Data:", err.Data)
		os.Exit(1)
	} else {
		fmt.Println(result)
	}
}

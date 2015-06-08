package main

import (
	"fmt"
	"github.com/THECALLR/sdk-go"
)

func main() {

	// initialize instance Thecallr
	thecallr.Setup("login", "password", nil)

	// an optional third parameter let you add options like proxy support
	// proxy must be in url standard format
	// http[s]://user:password@host:port
	// http[s]://host:port
	// http[s]://host

	// var config thecallr.Config
	// config.Proxy = "http://foo:bar@example.com:8080"
	// thecallr.Setup("login", "password", &config)
	

	// Example to send a SMS
	// 1. "call" method: each parameter of the method as an argument
	result, err := thecallr.Call("sms.send", "THECALLR", "+33123456789", "Hello, world", map[string]interface{}{
		"flash_message": false,
	})

	// error management
	if err != nil {
		fmt.Println("Code:", err.Code)
		fmt.Println("Message:", err.Msg)
		fmt.Println("Data:", err.Data)
	} else {
		fmt.Println(result)
	}

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

	// error management
	if err != nil {
		fmt.Println("Code:", err.Code)
		fmt.Println("Message:", err.Msg)
		fmt.Println("Data:", err.Data)
	} else {
		fmt.Println(result)
	}
}

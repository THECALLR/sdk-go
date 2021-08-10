// Package callr implements the CALLR API, using JSON-RPC 2.0. See https://www.callr.com/ and https://www.callr.com/docs/.

// This package may emit logs when errors occur when communicating with the API.
// The default logging function is log.Printf from the standard library. You can change
// the logging function with SetLogFunc.
//
// Usage
//
//    package main
//
//    import (
//        "context"
//        "fmt"
//        "os"
//
//        callr "github.com/THECALLR/sdk-go/v2"
//    )
//
//    func main() {
//        // use Basic Auth (not recommended)
//        // api := callr.NewWithBasicAuth("login", "password")
//
//        // or use Api Key Auth (recommended)
//        api := api.NewWithAPIKeyAuth("key")
//
//        // optional: set a proxy
//        // api.SetProxy("http://proxy:port")
//
//        // check for destination phone number parameter
//        if len(os.Args) < 2 {
//            fmt.Println("Please supply destination phone number!")
//            os.Exit(1)
//        }
//
//        // Example to send an SMS
//        result, err := api.Call(context.Background(), "sms.send", "SMS", os.Args[1], "Hello, world", nil)
//
//        // error management
//        if err != nil {
//            var jsonRpcError *callr.JSONRPCError
//            if errors.As(err, &jsonRpcError) {
//                fmt.Printf("API error: code:%d message:%s data:%v\n", jsonRpcError.Code, jsonRpcError.Message, jsonRpcError.Data)
//            } else {
//                fmt.Println("Transport error: ", err)
//            }
//            os.Exit(1)
//        }
//
//        fmt.Println(result)
//    }
package callr

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
)

// internal types

type jsonRCPRequest struct {
	ID      int64         `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type jsonRPCResponse struct {
	ID      int64         `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// API represents a connection to the CALLR API.
type API struct {
	urls   []string
	auth   string
	client *http.Client
}

// JSONRPCError is a JSON-RPC 2.0 error, returned by the API. It satisfies the native error interface.
type JSONRPCError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LogFunc func(string, ...interface{}) // Printf style

const (
	apiURL         = "https://api.callr.com/json-rpc/v1.1/"
	sdkVersion     = "2.0"
	jsonrpcVersion = "2.0"
	maxRetries     = 3 // on multiple URLs
)

var (
	defaultURLs = []string{apiURL}
	logFunc     = log.Printf
)

// NewWithBasicAuth returns an API object with Basic Authentication (not recommended). Use NewWithAPIKeyAuth auth instead.
func NewWithBasicAuth(login, password string) *API {
	return &API{
		urls:   defaultURLs,
		auth:   "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password)),
		client: &http.Client{},
	}
}

// NewWithAPIKeyAuth returns an API object with an API Key Authentication.
func NewWithAPIKeyAuth(key string) *API {
	return &API{
		urls:   defaultURLs,
		auth:   "Api-Key " + key,
		client: &http.Client{},
	}
}

// Error implements the error interface. Returns a string with the Code and Message properties.
func (e *JSONRPCError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// SetLogFunc can be used to change the default logger (log.Printf). Set to nil to disable package logging.
func SetLogFunc(fn LogFunc) error {
	if fn == nil {
		fn = func(string, ...interface{}) {
			// do nothing
		}
	}

	logFunc = fn
	return nil
}

// SetURL changes the URL for the API object
func (api *API) SetURL(url string) error {
	api.urls = []string{url}
	return nil
}

// SetURLs sets multiple URL for the API object. One URL is randomly selected when querying the API.
func (api *API) SetURLs(urls []string) error {
	if urls == nil {
		return errors.New("urls cannot be nil")
	}

	api.urls = urls
	return nil
}

// SetProxy sets a proxy URL to use
func (api *API) SetProxy(proxy string) error {
	url, err := url.Parse(proxy)

	if err != nil {
		return err
	}

	api.client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(url),
		},
	}

	return nil
}

// Call sends a JSON-RPC 2.0 request to the CALLR API, and returns either a result or an error.
// The error may be of type *JSONRPCError if the error comes from the API, or a native error otherwise.
func (api *API) Call(ctx context.Context, method string, params ...interface{}) (interface{}, error) {
	if params == nil {
		params = []interface{}{} // empty array instead of null
	}

	request := jsonRCPRequest{
		ID:      rand.Int63(),
		Method:  method,
		Params:  params,
		JSONRPC: jsonrpcVersion,
	}

	body, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	urls := make([]string, len(api.urls))
	copy(urls, api.urls)

	var lastError error

	for try := 0; try <= maxRetries; try++ {
		if len(urls) == 0 {
			return nil, lastError
		}

		randomIndex := rand.Intn(len(urls))
		url := urls[randomIndex]

		// slice index
		if randomIndex < len(urls)-1 {
			urls = append(urls[0:randomIndex], urls[randomIndex+1:]...)
		} else {
			urls = urls[:len(urls)-1]
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))

		if err != nil {
			return nil, err
		}

		req.Header.Add("Authorization", api.auth)
		req.Header.Add("Content-Type", "application/json-rpc; charset=utf-8")
		req.Header.Add("User-Agent", fmt.Sprintf("sdk=GO; sdk-version=%s; lang-version=%s; platform=%s",
			sdkVersion, runtime.Version(), runtime.GOOS))

		resp, err := api.client.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}

		if err != nil {
			lastError = err
			logFunc("[warning] url \"%s\" error: %s\n", url, err)
			// retry
			continue
		}

		var buf []byte

		if buf, err = io.ReadAll(resp.Body); err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			lastError = fmt.Errorf("response code: %d", resp.StatusCode)
			logFunc("[warning] url \"%s\" response code: %d\n", url, resp.StatusCode)
			// retry
			continue
		}

		if try > 0 {
			logFunc("[warning] successful at try: %d, on url: %s\n", try, url)
		}

		jsonResponse := jsonRPCResponse{}

		if err = json.Unmarshal(buf, &jsonResponse); err != nil {
			return nil, err
		}

		if jsonResponse.Error != nil {
			return nil, jsonResponse.Error
		}

		return jsonResponse.Result, nil
	}

	if lastError == nil {
		lastError = errors.New("unknown error")
	}

	return nil, lastError
}

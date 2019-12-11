// Package callr implements the CALLR API, using JSON-RPC 2.0. See https://www.callr.com/ and https://www.callr.com/docs/.
//
// Usage
//
//    package main
//
//    import (
//        "fmt"
//        "os"
//
//        callr "github.com/THECALLR/sdk-go"
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
//        // Example to send a SMS
//        result, err := api.Call("sms.send", "SMS", os.Args[1], "Hello, world", nil)
//
//        // error management
//        if err != nil {
//            if e, ok := err.(*callr.JSONRPCError); ok {
//                fmt.Printf("Remote error: code:%d message:%s data:%v\n", e.Code, e.Message, e.Data)
//            } else {
//                fmt.Println("Local error: ", err)
//            }
//            os.Exit(1)
//        }
//
//        fmt.Println(result)
//    }
package callr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"
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
	url    string
	auth   string
	client *http.Client
}

// JSONRPCError is a JSON-RPC 2.0 error, returned by the API. It satisfies the native error interface.
type JSONRPCError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	apiURL         = "https://api.callr.com/json-rpc/v1.1/"
	sdkVersion     = "1.0"
	jsonrpcVersion = "2.0"
)

// NewWithBasicAuth returns an API object with Basic Authentication (not recommended). Use NewWithAPIKeyAuth auth instead.
func NewWithBasicAuth(login, password string) *API {
	return &API{
		url:    apiURL,
		auth:   "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password)),
		client: &http.Client{},
	}
}

// NewWithAPIKeyAuth returns an API object with an API Key Authentication.
func NewWithAPIKeyAuth(key string) *API {
	return &API{
		url:    apiURL,
		auth:   "Api-Key " + key,
		client: &http.Client{},
	}
}

// Error implements the error interface. Returns a string with the Code and Message properties.
func (e *JSONRPCError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// SetURL changes the URL for the API object
func (api *API) SetURL(url string) error {
	api.url = url
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

// Call sends a JSON-RPC 2.0 request to the CALLR API, and returns either a result, or an error. The error may be of type JSONRPCError if the error comes from the API, or a native error if the error is local.
func (api *API) Call(method string, params ...interface{}) (interface{}, error) {
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

	req, err := http.NewRequest("POST", api.url, bytes.NewBuffer(body))

	req.Header.Add("Authorization", api.auth)
	req.Header.Add("Content-Type", "application/json-rpc; charset=utf-8")
	req.Header.Add("User-Agent", fmt.Sprintf("sdk=GO; sdk-version=%s; lang-version=%s; platform=%s",
		sdkVersion, runtime.Version(), runtime.GOOS))

	resp, err := api.client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	var buf []byte

	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

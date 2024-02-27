/*
Package callr implements the Callr API, using JSON-RPC 2.0. See https://www.callr.com/ and https://www.callr.com/docs/.

Requires Go 1.22+, using standard packages only.

Note:
This package may emit logs when errors occur when communicating with the API.
The default logging function is log.Printf from the standard library. You can change
the logging function with [SetLogFunc].

Usage:

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
*/
package callr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
)

// internal types

type jsonRCPRequest struct {
	ID      int64  `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type jsonRPCResponse struct {
	ID      int64           `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// API represents a connection to the Callr API.
type API struct {
	urls         []string
	auth         string
	client       *http.Client
	loginAsType  LoginAsType
	loginAsValue string
}

// JSONRPCError is a JSON-RPC 2.0 error, returned by the API. It satisfies the native error interface.
type JSONRPCError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// HTTPError is an HTTP error with a code and a message. It satisfies the native error interface.
type HTTPError struct {
	Code    int
	Message string
}

type LogFunc func(string, ...any) // Printf style
type LoginAsType string

const (
	apiURL         = "https://api.callr.com/json-rpc/v1.1/"
	sdkVersion     = "2.0"
	jsonrpcVersion = "2.0"
	maxRetries     = 3 // on multiple URLs

	LoginAsAccountID   LoginAsType = "account.id"
	LoginAsAccountRef  LoginAsType = "account.hash"
	LoginAsAccountHash LoginAsType = "account.hash"
	LoginAsUserID      LoginAsType = "user.id"
	LoginAsUserLogin   LoginAsType = "user.login"
)

var (
	defaultURLs       = []string{apiURL}
	logFunc           = log.Printf
	validLoginAsTypes = []LoginAsType{
		LoginAsAccountID,
		LoginAsAccountRef,
		LoginAsAccountHash,
		LoginAsUserID,
		LoginAsUserLogin,
	}
)

// NewWithAPIKeyAuth returns an [API] object with API Key Authentication.
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

// Error implements the error interface. Returns a string with the Code and Message properties.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// SetLogFunc can be used to change the default logger (log.Printf). Set to nil to disable package logging.
func SetLogFunc(fn LogFunc) {
	if fn == nil {
		fn = func(string, ...any) {
			// do nothing
		}
	}

	logFunc = fn
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

// SetLoginAs allows you to connect to the API as a sub account of yours, using different target types.
func (api *API) SetLoginAs(targetType LoginAsType, value string) error {
	if len(targetType) == 0 || len(value) == 0 {
		return errors.New("invalid login-as target type or value")
	}

	found := false

	for _, target := range validLoginAsTypes {
		if targetType == target {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid login-as target type: %s", targetType)
	}

	api.loginAsType = targetType
	api.loginAsValue = value

	return nil
}

// SetLoginAsSubAccountRef allows you to connect to the API as a sub-account of yours, using the sub-account "ref" field (sometimes called "hash").
func (api *API) SetLoginAsSubAccountRef(accountRef string) error {
	if len(accountRef) == 0 {
		return errors.New("invalid empty login-as account.ref")
	}

	api.loginAsType = LoginAsAccountRef
	api.loginAsValue = accountRef

	return nil
}

// SetLoginAsSubAccountLogin allows you to connect to the API as a sub-account of yours, using the sub-account user "login" field.
func (api *API) SetLoginAsSubAccountLogin(userLogin string) error {
	if len(userLogin) == 0 {
		return errors.New("invalid empty login-as account.ref")
	}

	api.loginAsType = LoginAsUserLogin
	api.loginAsValue = userLogin

	return nil
}

// ResetLoginAs removes the login-as configuration.
func (api *API) ResetLoginAs() {
	api.loginAsType = ""
	api.loginAsValue = ""
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

// Call sends a JSON-RPC 2.0 request to the Callr API, and returns either a result or an error.
// The error may be of type [*JSONRPCError] if the error comes from the API, of type [*HTTPError] if the error comes from the HTTP transport,
// or a native error otherwise.
func (api *API) Call(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	if params == nil {
		params = []any{} // empty array instead of null
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

		if len(api.loginAsType) != 0 {
			req.Header.Add("Callr-Login-As", fmt.Sprintf("%s %s", api.loginAsType, api.loginAsValue))
		}

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
			lastError = newHTTPError(resp)
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

func newHTTPError(resp *http.Response) error {
	status := resp.Status
	code := strconv.Itoa(resp.StatusCode)

	// remove the status code from the message
	status = strings.TrimPrefix(status, code)

	return &HTTPError{
		Code:    resp.StatusCode,
		Message: strings.TrimLeft(status, " "),
	}
}

/**
* CALLR webservice communication library
**/

package callr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"runtime"
	"fmt"
)

type Auth interface {
	Apply(request Request)
	LogAs(as_type string, as_target string) Error
}

type LoginPasswordAuth struct {
	login	 string
	password string
	logAs	 string
}

func LoginPassword(login string, password string) LoginPasswordAuth {
	return LoginPasswordAuth{login, password}
}

func (LoginPasswordAuth *auth) Apply(request Request) {
	encoded = base64.StdEncoding.EncodeToString([]byte(auth.login + ":" + auth.password))

	request.Header.Add("Authorization", "Basic "+encoded)

	if auth.logAs {
		request.Header.add("CALLR-Login-As", auth.logAs)
	}
}

func (LoginPasswordAuth *auth) LogAs(as_type string, as_target string) {
	if as_type == "" && as_target == "" {
		auth.logAs = ""
		return
	}

	switch strings.toLower(as_type) {
	case "user":
		as_type = "User.login"
		break
	
	case "account":
		as_type = "Account.hash"
		break
	
	default:
		return NewError("LOGIN_AS_WRONG_TYPE", 2, nil)
	}
 
	auth.logAs = fmt.Sprintf("%s %s", as_type, as_target)
}

type Callr struct {
	Auth	 Auth
	ApiUrl   string
	Config   *Config
}

type Config struct {
	Proxy string
}

type Json struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Response struct {
	Body string
	Json map[string]interface{}
}

type Error struct {
	Msg  string
	Code int
	Data map[string]interface{}
}

var TC Callr
const SDK_VERSION = "0.1"

/*******************************************************************************
*** Callr Methods
*******************************************************************************/

/**
* Convert Credentials to base64 using "Login:Password" format
**/
func (t Callr) Base64() string {
	return base64.StdEncoding.EncodeToString([]byte(t.Login + ":" + t.Password))
}

/*******************************************************************************
*** Json Methods
*******************************************************************************/

/**
* Convert object (struct) to JSON string
**/
func (obj *Json) ToString() (data []byte) {
	data, _ = json.Marshal(obj)
	return
}

/**
* Format params to object
**/
func (obj *Json) make(method string, params []interface{}, id []int) *Json {
	obj.Jsonrpc = "2.0"
	obj.Method = method
	obj.Params = params
	if len(id) != 0 {
		obj.Id = id[0]
	} else {
		obj.Id = 100 + rand.Intn(999-100)
	}
	return obj
}

/*******************************************************************************
*** SDK Functions
*******************************************************************************/

func init() {
	rand.Seed(time.Now().UnixNano()) // Random seed generator
	TC.ApiUrl = "https://api.callr.com/json-rpc/v1.1/"
}

/**
* Initialisation
* @param string Login
* @param string Password
**/
func Setup(login, password string, config *Config) {
	if login != "" && password != "" {
		TC.Auth = LoginPasswordAuth(login, password)
	}

	TC.Config = config
}

/**
* Send a request to CALLR webservice
**/
func Call(args ...interface{}) (*Response, *Error) {
	return Send(args[0].(string), args[1:])
}

/**
* Send a request to CALLR webservice
**/
func Send(method string, params []interface{}, id ...int) (*Response, *Error) {
	if TC.Auth == nil {
		return nil, NewError("CREDENTIALS_NOT_SET", -1, nil)
	}

	// create object for json encode
	object := new(Json).make(method, params, id)

	// Create Request
	// Proxy support
	var client *http.Client
	if TC.Config != nil && len(TC.Config.Proxy) > 0 {
		proxyUrl, err := url.Parse(TC.Config.Proxy)
		if err != nil {
			return nil, NewError("PROXY_PARSE_URL_ERROR", 1, nil)
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		client = &http.Client{}
	}
	req, err := http.NewRequest("POST", TC.ApiUrl, bytes.NewBuffer(object.ToString()))

	req.Header.Add("Content-Type", "application/json-rpc; charset=utf-8")
	req.Header.Add("User-Agent", fmt.Sprintf("sdk=GO; sdk-version=%s; lang-version=%s; platform=%s", SDK_VERSION, runtime.Version(), runtime.GOOS))

	TC.Auth.Apply(req)

	resp, err := client.Do(req)

	// Error management
	if err != nil {
		return nil, NewError(err.Error(), -1, nil)
	} else if resp.StatusCode != 200 {
		return nil, NewError("HTTP_CODE_ERROR", resp.StatusCode, nil)
	}

	return ParseResponse(resp)
}

/**
* Response analysis
**/
func ParseResponse(r *http.Response) (resp *Response, error *Error) {
	defer r.Body.Close()

	content, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, NewError(err.Error(), -1, nil)
	}

	resp = new(Response)
	resp.Body = string(content)
	err = json.Unmarshal(content, &resp.Json)

	// Error management
	error = nil
	if err != nil {
		return nil, NewError("INVALID_RESPONSE", -1, map[string]interface{}{"response": resp})
	} else if error, ok := resp.Json["error"]; ok {
		return nil, NewError(error.(map[string]interface{})["message"].(string), int(error.(map[string]interface{})["code"].(float64)), nil)
	}

	return
}

func NewError(s string, code int, data map[string]interface{}) (err *Error) {
	err = new(Error)
	err.Msg = s
	err.Code = code
	err.Data = data
	return err
}

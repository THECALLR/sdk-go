/**
* THECALLR webservice communication library
**/

package thecallr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Thecallr struct {
	Login    string
	Password string
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

var TC Thecallr

/*******************************************************************************
*** Thecallr Methods
*******************************************************************************/

/**
* Convert Credentials to base64 using "Login:Password" format
**/
func (t Thecallr) Base64() string {
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
	TC.ApiUrl = "https://api.thecallr.com"
}

/**
* Initialisation
* @param string Login
* @param string Password
**/
func Setup(login, password string, config *Config) {
	TC.Login = login
	TC.Password = password
	TC.Config = config
}

/**
* Send a request to THECALLR webservice
**/
func Call(args ...interface{}) (*Response, *Error) {
	return Send(args[0].(string), args[1:])
}

/**
* Send a request to THECALLR webservice
**/
func Send(method string, params []interface{}, id ...int) (*Response, *Error) {
	auth := CheckAuth()
	if auth != nil {
		return nil, auth
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

	req.Header.Add("Authorization", "Basic "+TC.Base64())
	req.Header.Add("Content-Type", "application/json-rpc; charset=utf-8")

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

func CheckAuth() *Error {
	if len(TC.Login) == 0 || len(TC.Password) == 0 {
		return NewError("CREDENTIALS_NOT_SET", -1, nil)
	}
	return nil
}

func NewError(s string, code int, data map[string]interface{}) (err *Error) {
	err = new(Error)
	err.Msg = s
	err.Code = code
	err.Data = data
	return err
}

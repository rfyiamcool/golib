package requests

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

// Supported http methods
const (
	MethodGet    string = "GET"
	MethodPost   string = "POST"
	MethodPut    string = "PUT"
	MethodPatch  string = "PATCH"
	MethodDelete string = "DELETE"
)

var (
	defaultClient *http.Client

	debugMode = false
)

func init() {
	var dials uint64
	tr := &http.Transport{}
	tr.Dial = func(network, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(network, addr, time.Second*10)
		if err == nil && debugMode {
			log.Printf("rpc: dial new connection to %s， dials: %d \n", addr, atomic.AddUint64(&dials, 1)-1)
		}
		return c, err
	}
	defaultClient = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
}

// Request request info
type Request struct {
	User     string
	Password string
	Method   string
	URL      string
	Header   map[string]string
	Params   map[string]string
	Cookies  map[string]string
	Body     []byte
}

func NewRequest(url, method string, header map[string]string, params map[string]string, data []byte) *Request {
	return &Request{Method: method, URL: url, Header: header, Params: params, Body: data}
}

func (req *Request) SetBasicAuth(user, password string) {
	req.User = user
	req.Password = password
}

type Response struct {
	StatusCode int
	Header     map[string][]string
	Body       []byte
}

func AddParameters(baseURL string, queryParams map[string]string) string {
	baseURL += "?"
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return baseURL + params.Encode()
}

func BuildHTTPRequest(request *Request) (*http.Request, error) {
	var (
		err     error
		httpReq *http.Request
		body    io.Reader
	)

	// handle parameters
	if request.Method == MethodPost {
		args := url.Values{}
		for k, v := range request.Params {
			args.Set(k, v)
		}

		body = strings.NewReader(args.Encode())
	} else {
		if len(request.Params) > 0 {
			request.URL = AddParameters(request.URL, request.Params)
		}

		body = bytes.NewReader(request.Body)
	}

	// build http request
	httpReq, err = http.NewRequest(request.Method, request.URL, body)
	if err != nil {
		return nil, err
	}

	// set basic auth
	if request.User != "" && request.Password != "" {
		httpReq.SetBasicAuth(request.User, request.Password)
	}

	// build http header
	for k, v := range request.Header {
		httpReq.Header.Set(k, v)
	}

	// default type
	_, ok := request.Header["Content-Type"]
	if len(request.Body) > 0 && !ok {
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return httpReq, nil
}

func buildResponse(res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	r := &Response{
		StatusCode: res.StatusCode,
		Body:       body,
		Header:     res.Header,
	}
	return r, nil
}

// Send send http request
func Send(request *Request) (*Response, error) {
	var start = time.Now()

	// build http request
	httpReq, err := BuildHTTPRequest(request)
	if err != nil {
		return nil, err
	}

	// send http request
	rsp, err := defaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer func() {
		io.Copy(ioutil.Discard, rsp.Body)
		// close http response
		rsp.Body.Close()

		if debugMode {
			log.Printf("call rpc [%s] %s in %v \n", httpReq.Method, httpReq.URL, time.Since(start))
		}
	}()

	// build response
	r, err := buildResponse(rsp)
	return r, err
}

func PostBody(url string, header map[string]string, data []byte) (*Response, error) {
	req := NewRequest(url, MethodGet, header, nil, data)
	req.Header["Content-Type"] = "application/json"
	resp, err := Send(req)
	return resp, err
}

func Post(api string, header map[string]string, params map[string]string) (*Response, error) {
	req := NewRequest(api, MethodGet, header, params, nil)
	resp, err := Send(req)
	return resp, err
}

func Get(url string, header map[string]string, params map[string]string) (*Response, error) {
	req := NewRequest(url, MethodGet, header, params, nil)
	resp, err := Send(req)
	return resp, err
}

func EncodeURL(host string, format string, args ...interface{}) string {
	var u url.URL
	u.Scheme = "http"
	u.Host = host
	u.Path = fmt.Sprintf(format, args...)
	return u.String()
}

func SetDebug(b bool) {
	debugMode = b
}

func SetTimeout(d time.Duration) {
	defaultClient.Timeout = d
}

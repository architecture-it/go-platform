package httpExtension

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/architecture-it/go-platform/log"
)

type Params map[string]string

// Default empty
var DefaultHeader http.Header

// The timeout includes connection time, any
// redirects, and reading the response body. The timer remains
// running after Get, Head, Post, or Do return and will
// interrupt reading of the Response.Body.
//
// Default 10 seconds. A Timeout of zero means no timeout.
//
// Read more in http package.
var Timeout time.Duration

func init() {
	WithTimeout(10)
	DefaultHeader = http.Header{}
}

// Send a GET request to the specified URL
var Get = func(requestUrl string, params Params, header http.Header) (*Response, error) {
	return get(requestUrl, params, header, Timeout)
}

// Send a GET request to the specified URL as an asynchronous operation
var GetAsync = func(requestUrl string, params map[string]string, header http.Header) <-chan *Response {
	c := make(chan *Response)
	timeout := Timeout

	go func() {
		res, _ := get(requestUrl, params, header, timeout)
		c <- res
		close(c)
	}()

	return c
}

// Send a POST request to the specified Url
var Post = func(requestUrl string, body []byte, header http.Header) (*Response, error) {
	return post(requestUrl, body, header, Timeout)
}

// Send a POST request to the specified Url as an asynchronous operation
var PostAsync = func(requestUrl string, body []byte, header http.Header) <-chan *Response {
	c := make(chan *Response)
	timeout := Timeout

	go func() {
		res, _ := post(requestUrl, body, header, timeout)
		c <- res
		close(c)
	}()

	return c
}

// Set timeout expressed in seconds
func WithTimeout(seconds int64) {
	s := time.Duration(seconds)
	Timeout = s * time.Second
}

func get(requestUrl string, params Params, header http.Header, timeout time.Duration) (*Response, error) {
	u, err := url.Parse(requestUrl)
	if err != nil {
		log.Logger.Error("Couldn't parse url " + requestUrl + ". Error: " + err.Error())
		return nil, err
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return doRequest(http.MethodGet, u.String(), nil, header, timeout)
}

func post(requestUrl string, body []byte, header http.Header, timeout time.Duration) (*Response, error) {
	return doRequest(http.MethodPost, requestUrl, bytes.NewBuffer(body), header, timeout)
}

func doRequest(method, requestUrl string, body io.Reader, header http.Header, timeout time.Duration) (*Response, error) {
	req, err := http.NewRequest(method, requestUrl, body)
	if err != nil {
		log.Logger.Error("Failed to create request. Error: " + err.Error())
		return nil, err
	}

	req.Header = DefaultHeader
	addHeader(req, header)

	httpClient := &http.Client{Timeout: timeout}
	res, err := httpClient.Do(req)

	if err != nil {
		log.Logger.Error(method + " " + requestUrl + ". Error: " + err.Error())
		return nil, err
	}

	if res == nil {
		log.Logger.Error(method + " " + requestUrl + ". Nil response")
		return nil, errors.New("nil response")
	}

	response := Response{Response: res}
	return &response, nil
}

func addHeader(req *http.Request, header http.Header) {
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

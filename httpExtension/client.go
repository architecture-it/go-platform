package httpExtension

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/architecture-it/go-platform/log"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
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
var Get = func(requestUrl string, params Params, header http.Header, ctx context.Context) (*Response, error) {
	return get(requestUrl, params, header, Timeout, ctx)
}

// Send a GET request to the specified URL as an asynchronous operation
var GetAsync = func(requestUrl string, params Params, header http.Header, ctx context.Context) <-chan *Response {
	c := make(chan *Response)
	timeout := Timeout

	go func() {
		res, _ := get(requestUrl, params, header, timeout, ctx)
		c <- res
		close(c)
	}()

	return c
}

// Send a POST request to the specified Url
var Post = func(requestUrl string, body []byte, header http.Header, ctx context.Context) (*Response, error) {
	return post(requestUrl, body, header, Timeout, ctx)
}

// Send a POST request to the specified Url as an asynchronous operation
var PostAsync = func(requestUrl string, body []byte, header http.Header, ctx context.Context) <-chan *Response {
	c := make(chan *Response)
	timeout := Timeout

	go func() {
		res, _ := post(requestUrl, body, header, timeout, ctx)
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

func get(requestUrl string, params Params, header http.Header, timeout time.Duration, ctx context.Context) (*Response, error) {
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

	return doRequest(http.MethodGet, u.String(), nil, header, timeout, ctx)
}

func post(requestUrl string, body []byte, header http.Header, timeout time.Duration, ctx context.Context) (*Response, error) {
	return doRequest(http.MethodPost, requestUrl, bytes.NewBuffer(body), header, timeout, ctx)
}

func doRequest(method, requestUrl string, body io.Reader, header http.Header, timeout time.Duration, ctx context.Context) (*Response, error) {

	req, err := http.NewRequestWithContext(ctx, method, requestUrl, body)
	if err != nil {
		log.Logger.Error("Failed to create request. Error: " + err.Error())
		return nil, err
	}

	addHeader(req, DefaultHeader)
	addHeader(req, header)

	httpClient := &http.Client{Timeout: timeout}
	var tracingClient = apmhttp.WrapClient(httpClient)
	res, err := tracingClient.Do(req)

	if err != nil {
		log.Logger.Error(method + " " + requestUrl + ". Error: " + err.Error())
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	if res == nil {
		log.Logger.Error(method + " " + requestUrl + ". Nil response")
		apm.CaptureError(ctx, err).Send()
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

/*
func main() {
	var myHandler http.Handler = ...
	tracedHandler := apmhttp.Wrap(myHandler)
}

The apmhttp handler will recover panics and send them to Elastic APM.

Package apmhttp also provides functions for instrumenting an http.Client or http.RoundTripper such that outgoing requests are traced as spans, if the request context includes a transaction. When performing the request, the enclosing context should be propagated by using http.Request.WithContext, or a helper, such as those provided by https://golang.org/x/net/context/ctxhttp.

Client spans are not ended until the response body is fully consumed or closed. If you fail to do either, the span will not be sent. Always close the response body to ensure HTTP connections can be reused; see func (*Client) Do.


var tracingClient = apmhttp.WrapClient(http.DefaultClient)

func serverHandler(w http.ResponseWriter, req *http.Request) {
	// Propagate the transaction context contained in req.Context().
	resp, err := ctxhttp.Get(req.Context(), tracingClient, "http://backend.local/foo")
	if err != nil {
		apm.CaptureError(req.Context(), err).Send()
		http.Error(w, "failed to query backend", 500)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	...
}

func main() {
	http.ListenAndServe(":8080", apmhttp.Wrap(http.HandlerFunc(serverHandler)))
}*/

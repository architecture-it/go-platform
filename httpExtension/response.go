package httpExtension

import (
	"errors"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
}

// true if StatusCode was in the range 100-199; otherwise, false.
func (r *Response) IsInformationalStatusCode() bool {
	return r.StatusCode >= 100 && r.StatusCode <= 199
}

// true if StatusCode was in the range 200-299; otherwise, false.
func (r *Response) IsSuccessStatusCode() bool {
	return r.StatusCode >= 200 && r.StatusCode <= 299
}

// true if StatusCode was in the range 400-499; otherwise, false.
func (r *Response) IsClientErrorStatusCode() bool {
	return r.StatusCode >= 400 && r.StatusCode <= 499
}

// true if StatusCode was in the range 500-599; otherwise, false.
func (r *Response) IsServerErrorStatusCode() bool {
	return r.StatusCode >= 500 && r.StatusCode <= 599
}

// true if StatusCode was in the parameter list; otherwise, false.
func (r *Response) IsStatusCodeIn(status []int) bool {
	for _, status := range status {
		if r.StatusCode == status {
			return true
		}
	}

	return false
}

// read all body and close it
func (r *Response) ReadAndClose() ([]byte, error) {
	if r.Body == nil {
		return []byte{}, errors.New("body can't be nil")
	}

	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

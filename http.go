/*
Package http seeks to effectively take the place of curl for Go
rapid-applications development. It provides very high-level HTTP
functions to simplify web requests. The standard net/http library is
robust and complete, but requires far more lines of code to implement
anything quickly.

Leveraging Generics

By using generics we can pass different receptacle structures to receive the HTTP responses.
*/
package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// TimeOut is a package global timeout for any of the high-level https
// query functions in this package. The default value is 60 seconds.
var TimeOut int = 60

// MaxRedirects is the maximum number of redirects any request will
// attempt before returning an error. The default value is 5.
var MaxRedirects = 5

// Get is much like http.Get except that it unmarshals into the
// specified struct (which may already contain populated data fields).
// Get also observes the package global http.TimeOut and will
// automatically redirect when a redirect response is received. In fact,
// errors are returned for any status code other than anything in the
// 200 range (including after a successful redirect).
func Get[T any](url string, data *T) error {
	var err error

	// build request with no body
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	//  upgrade request to with context and TimeOut
	dur := time.Duration(time.Second * time.Duration(TimeOut))
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()
	req = req.WithContext(ctx)

	// do the request and check status code, if redirect do it

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// TODO redirect 300 status responses

	// read all the body of the response and unmarshal it
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, data)
}

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

// Client provides a way to change the default HTTP client for
// any further package HTTP request function calls. By default, it is
// set to http.DefaultClient. This is particularly useful when creating
// mockups and other testing.
var Client = http.DefaultClient

// Request passes the requested method with the given URL and input data
// values to the HTTP Client and unmarshals the response into the data
// struct passed by pointer (out, which may already contain populated
// data fields). Request also observes the package global http.TimeOut
// Any status code other than 200 returns an error. Also see
// Get, Post, Put, Patch, and Delete.
func Request[T any](method, url string, in url.Values, out *T) error {
	var err error
	var inreader io.Reader
	var inlength string

	// encode any input data
	if in != nil {
		encoded := in.Encode()
		inreader = strings.NewReader(encoded)
		inlength = strconv.Itoa(len(encoded))
	}

	// build request with no body
	req, err := http.NewRequest(method, url, inreader)
	if in != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", inlength)
	}

	if err != nil {
		return err
	}

	//  upgrade request to with context and TimeOut
	dur := time.Duration(time.Second * time.Duration(TimeOut))
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()
	req = req.WithContext(ctx)

	// do the request and check status code
	res, err := Client.Do(req)
	if err != nil {
		return err
	}

	// read all the body of the response and unmarshal it
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, out)
}

// Get sends a GET Request.
func Get[T any](url string, out *T) error {
	return Request(`GET`, url, nil, out)
}

// Post sends a POST Request.
func Post[T any](url string, in url.Values, out *T) error {
	return Request(`POST`, url, in, out)
}

// Put sends a POST Request.
func Put[T any](url string, in url.Values, out *T) error {
	return Request(`PUT`, url, in, out)
}

// Patch sends a PATCH Request.
func Patch[T any](url string, in url.Values, out *T) error {
	return Request(`PATCH`, url, in, out)
}

// Delete sends a DELETE Request.
func Delete[T any](url string, out *T) error {
	return Request(`DELETE`, url, nil, out)
}

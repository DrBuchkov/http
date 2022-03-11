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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// Req encapsulates a single HTTP GET Request in a way that can be
// combined with others Requestions through Pipe.
type Req struct {
	Method string
	URL    string
	Data   url.Values
}

// GET is shorthand for Req(`GET`, ...). Data url.Values are passed as
// a query string in the URL.
type GET struct {
	URL  string
	Data url.Values
}

// POST is shorthand for Req(`POST`, ...). Data url.Values are passed as
// encoded form data in the body.
type POST struct {
	URL  string
	Data url.Values
}

// PUT is shorthand for Req(`PUT`, ...). Data url.Values are passed as
// encoded form data in the body.
type PUT struct {
	URL  string
	Data url.Values
}

// PATCH is shorthand for Req(`PATCH`, ...). Data url.Values are passed
// encoded as form data in the body.
type PATCH struct {
	URL  string
	Data url.Values
}

// DELETE is shorthand for Req(`DELETE`, ...). Data url.Values are
// passed as a query string in the URL.
type DELETE struct {
	URL  string
	Data url.Values
}

// Request passes the requested method with the given URL and input data
// values to the HTTP Client and unmarshals the response into the data
// struct passed by pointer (out, which may already contain populated
// data fields). Request also observes the package global http.TimeOut
// Any status code other than 200 returns an error. Also see
// Get, Post, Put, Patch, and Delete.
func Request[T any](method, uri string, in url.Values, out *T) error {
	var err error
	var req *http.Request

	// encode any input data
	switch method {
	case "GET", "DELETE":
		req, err = http.NewRequest(method, uri, nil)
		if in != nil {
			q := req.URL.Query()
			for k, values := range in {
				for _, value := range values {
					q.Add(k, value)
				}
			}
			req.URL.RawQuery = q.Encode()
		}
		break
	case "POST", "PUT", "PATCH":
		var inreader *strings.Reader = nil
		var inlength string
		if in != nil {
			encoded := in.Encode()
			inreader = strings.NewReader(encoded)
			inlength = strconv.Itoa(len(encoded))
		}

		req, err = http.NewRequest(method, uri, inreader)
		if err != nil {
			return err
		}
		if in != nil {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("Content-Length", inlength)
		}
		break
	}

	//  upgrade request to with context and TimeOut
	dur := time.Duration(time.Second * time.Duration(TimeOut))
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if req == nil {
		fmt.Println("Hello World")
	}

	req = req.WithContext(ctx)

	// do the request and check status code
	res, err := Client.Do(req)
	if err != nil {
		return err
	}

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return fmt.Errorf(res.Status)
	}

	// read all the body of the response and unmarshal it
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, out)
}

// Get sends a GET Request.
func Get[T any](url string, in url.Values, out *T) error {
	return Request(`GET`, url, in, out)
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

// ReqRecipe is a "bottled" HTTP request, that can be used with http.Pipe.
type ReqRecipe[T any] func(out *T) error

func GetRecipe[T any](url string, in url.Values) ReqRecipe[T] {
	return ReqRecipe[T](func(out *T) error {
		return Get(url, in, out)
	})
}

func PostRecipe[T any](url string, in url.Values) ReqRecipe[T] {
	return ReqRecipe[T](func(out *T) error {
		return Post(url, in, out)
	})
}
func PutRecipe[T any](url string, in url.Values) ReqRecipe[T] {
	return ReqRecipe[T](func(out *T) error {
		return Put(url, in, out)
	})
}
func PatchRecipe[T any](url string, in url.Values) ReqRecipe[T] {
	return ReqRecipe[T](func(out *T) error {
		return Patch(url, in, out)
	})
}
func DeleteRecipe[T any](url string) ReqRecipe[T] {
	return ReqRecipe[T](func(out *T) error {
		return Delete(url, out)
	})
}

// Pipe Example of different flavor using higher order functions. In my experience code that tries
// to encapsulate "processes" using data structs ends up trying implementing an "interpreter" that blows up in complexity
// for every capability that one would want to support. Higher-order functions on the other hand are easily constructable,
// composable, and can leverage any kind of logic.
func Pipe[T any](data *T, recipes ...ReqRecipe[T]) error {
	for _, recipe := range recipes {
		if err := recipe(data); err != nil {
			return err
		}
	}
	return nil
}

// Pipe makes multiple HTTP requests in succession sending the same data
// object to all of them for marshaling with the results of the
// requests. This is useful for chaining multiple service or REST API
// requests together. It also allows chains of requests to be saved and
// added to registries for repeated and composition with other data flow
// pipelines.
//func Pipe[T any](data *T, reqs ...any) error {
//	for _, req := range reqs {
//		if req, isslice := req.([]Req); isslice {
//			for _, r := range req {
//				if err := Pipe(data, r); err != nil {
//					return err
//				}
//			}
//		}
//		switch v := req.(type) {
//		case GET:
//			if err := Get(v.URL, v.Data, data); err != nil {
//				return err
//			}
//		case POST:
//			if err := Post(v.URL, v.Data, data); err != nil {
//				return err
//			}
//		case PATCH:
//			if err := Patch(v.URL, v.Data, data); err != nil {
//				return err
//			}
//		case PUT:
//			if err := Put(v.URL, v.Data, data); err != nil {
//				return err
//			}
//		case DELETE:
//			if err := Delete(v.URL, data); err != nil {
//				return err
//			}
//		case Req:
//			switch v.Method {
//			case `GET`, `POST`, `PUT`, `PATCH`:
//				if err := Request(v.Method, v.URL, v.Data, data); err != nil {
//					return err
//				}
//			case `DELETE`:
//				if err := Delete(v.URL, data); err != nil {
//					return err
//				}
//			default:
//				return fmt.Errorf(`unsupported request method: %v`, v.Method)
//			}
//		default:
//			return fmt.Errorf(`unsupported request type: %T`, v)
//		}
//	}
//	return nil
//}

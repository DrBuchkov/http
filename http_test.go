package http_test

import (
	"fmt"
	_http "net/http"
	ht "net/http/httptest"

	"github.com/rwxrob/http"
	"github.com/rwxrob/json"
)

func ExampleGet() {

	// setup mock web service
	handler := _http.HandlerFunc(
		func(w _http.ResponseWriter, r *_http.Request) {
			fmt.Fprintf(w, `{"word":"hello","name":"Rob"}`)
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	// change global timeout
	http.TimeOut = 60

	// create the struct type matching the REST query JSON
	type Greeting struct {
		Word      string `json:"word"`
		Name      string `json:"name"`
		Untouched string `json:"untouched"`
	}

	// create empty greeting and fill it with request response
	// and capture the extra net/http response as well
	greet := Greeting{"OVERWRITEWORD", "NOT A NAME", "SAME"}
	if err := http.Get(svr.URL, &greet); err != nil {
		fmt.Println(err)
	}

	// "cast" to json.Object so we can use JSON and confirm parsing
	json.Object[Greeting]{greet}.Print()

	// Output:
	// {"word":"hello","name":"Rob","untouched":"SAME"}
}

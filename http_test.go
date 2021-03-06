package http_test

import (
	"fmt"
	_http "net/http"
	ht "net/http/httptest"
	"net/url"

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
	if err := http.Get(svr.URL, nil, &greet); err != nil {
		fmt.Println(err)
	}

	// "cast" to json.Object so we can use JSON and confirm parsing
	json.Object[Greeting]{greet}.Print()

	// Output:
	// {"word":"hello","name":"Rob","untouched":"SAME"}
}

func ExamplePost() {

	// setup mock web service
	handler := _http.HandlerFunc(
		func(w _http.ResponseWriter, r *_http.Request) {
			r.ParseForm()
			if len(r.Form["name"]) > 0 {
				out := fmt.Sprintf(`{"word":"hello","name":%q}`, r.Form["name"][0])
				fmt.Fprintf(w, out)
			}
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	// create the struct type matching the REST query JSON
	type Greeting struct {
		Word      string `json:"word"`
		Name      string `json:"name"`
		Untouched string `json:"untouched"`
	}

	data := url.Values{}
	data.Set("name", "Roberto")

	// create empty greeting and fill it with request response
	// and capture the extra net/http response as well
	greet := Greeting{"OVERWRITEWORD", "NOT A NAME", "SAME"}
	if err := http.Post(svr.URL, data, &greet); err != nil {
		fmt.Println(err)
	}

	// "cast" to json.Object so we can use JSON and confirm parsing
	json.Object[Greeting]{greet}.Print()

	// Output:
	// {"word":"hello","name":"Roberto","untouched":"SAME"}
}

func ExamplePut() {

	// setup mock web service
	handler := _http.HandlerFunc(
		func(w _http.ResponseWriter, r *_http.Request) {
			r.ParseForm()
			if len(r.Form["name"]) > 0 {
				out := fmt.Sprintf(`{"word":"hello","name":%q}`, r.Form["name"][0])
				fmt.Fprintf(w, out)
			}
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	// create the struct type matching the REST query JSON
	type Greeting struct {
		Word      string `json:"word"`
		Name      string `json:"name"`
		Untouched string `json:"untouched"`
	}

	data := url.Values{}
	data.Set("name", "Roberto")

	// create empty greeting and fill it with request response
	// and capture the extra net/http response as well
	greet := Greeting{"OVERWRITEWORD", "NOT A NAME", "SAME"}
	if err := http.Put(svr.URL, data, &greet); err != nil {
		fmt.Println(err)
	}

	// "cast" to json.Object so we can use JSON and confirm parsing
	json.Object[Greeting]{greet}.Print()

	// Output:
	// {"word":"hello","name":"Roberto","untouched":"SAME"}
}

func ExamplePatch() {

	// setup mock web service
	handler := _http.HandlerFunc(
		func(w _http.ResponseWriter, r *_http.Request) {
			r.ParseForm()
			if len(r.Form["name"]) > 0 {
				out := fmt.Sprintf(`{"word":"hello","name":%q}`, r.Form["name"][0])
				fmt.Fprintf(w, out)
			}
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	// create the struct type matching the REST query JSON
	type Greeting struct {
		Word      string `json:"word"`
		Name      string `json:"name"`
		Untouched string `json:"untouched"`
	}

	data := url.Values{}
	data.Set("name", "Roberto")

	// create empty greeting and fill it with request response
	// and capture the extra net/http response as well
	greet := Greeting{"OVERWRITEWORD", "NOT A NAME", "SAME"}
	if err := http.Patch(svr.URL, data, &greet); err != nil {
		fmt.Println(err)
	}

	// "cast" to json.Object so we can use JSON and confirm parsing
	json.Object[Greeting]{greet}.Print()

	// Output:
	// {"word":"hello","name":"Roberto","untouched":"SAME"}
}

func ExampleDelete() {

	// setup mock web service
	handler := _http.HandlerFunc(
		func(w _http.ResponseWriter, r *_http.Request) {
			fmt.Fprintf(w, `{"word":"hello","name":"Rob"}`)
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	// create the struct type matching the REST query JSON
	type Greeting struct {
		Word      string `json:"word"`
		Name      string `json:"name"`
		Untouched string `json:"untouched"`
	}

	data := url.Values{}
	data.Set("name", "Roberto")

	// create empty greeting and fill it with request response
	// and capture the extra net/http response as well
	greet := Greeting{"OVERWRITEWORD", "NOT A NAME", "SAME"}
	if err := http.Delete(svr.URL, &greet); err != nil {
		fmt.Println(err)
	}

	// "cast" to json.Object so we can use JSON and confirm parsing
	json.Object[Greeting]{greet}.Print()

	// Output:
	// {"word":"hello","name":"Rob","untouched":"SAME"}
}

func ExamplePipe() {

	// setup mock web service
	handler := _http.HandlerFunc(
		func(w _http.ResponseWriter, r *_http.Request) {
			switch r.URL.String() {
			case "/greet":
				fmt.Fprintf(w, `{"greeting":"hello"}`)
			case "/howru":
				fmt.Fprintf(w, `{"word":"hello","name":"Rob"}`)
			case "/bye":
				fmt.Fprintf(w, `{}`)
			default:
				w.WriteHeader(400)
			}
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	greet := http.GET{svr.URL + "/greet", nil}
	howru := http.POST{svr.URL + "/howru", url.Values{"name": []string{"Rob"}}}
	bye := http.GET{svr.URL + "/bye", nil}

	//encounter := []http.Req{greet, howru, bye}

	type Data struct {
		Name     string
		Greeting string
	}
	data := Data{}

	err := http.Pipe(&data, greet, howru, bye)
	if err != nil {
		fmt.Println(err)
	}

	json.Object[Data]{data}.Print()

	// Output:
	// {"Name":"Rob","Greeting":"hello"}
}

func ExampleGet_status() {

	// setup mock web service
	handler := _http.HandlerFunc(
		func(w _http.ResponseWriter, r *_http.Request) {
			w.WriteHeader(400)
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	type Data struct {
		Name     string
		Greeting string
	}
	data := Data{}

	if err := http.Get(svr.URL, nil, &data); err != nil {
		fmt.Println(err)
	}

	json.Object[Data]{data}.Print()

	// Output:
	// 400 Bad Request
	// {"Name":"","Greeting":""}
}

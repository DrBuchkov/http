package http_test

import (
	"fmt"
	_http "net/http"
	ht "net/http/httptest"
	"net/url"
	"strconv"

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

	type City struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		State       string `json:"state"`
		Population  int    `json:"population"`
		Temperature int    `json:"temperature"`
	}

	// Mocked "database" of cities
	cities := []City{
		{
			Id:    1,
			Name:  "New York",
			State: "New York",
		},
		{
			Id:    2,
			Name:  "Los Angeles",
			State: "California",
		},
		{
			Id:    3,
			Name:  "Chicago",
			State: "Illinois",
		},
	}
	// Mock handler that finds a city object given a name.
	cityServiceHandler := func(w _http.ResponseWriter, r *_http.Request) {

		r.ParseForm()
		if len(r.Form["name"]) > 0 {
			name := r.Form["name"][0]
			// Doesn't handle the case where no city is found but who cares
			var city City
			for _, c := range cities {
				if c.Name == name {
					city = c
				}
			}
			if city.Id != 0 {
				_, err := fmt.Fprintf(w, json.Object[City]{city}.String())
				if err != nil {
					panic(err)
				}
			}

		}

	}
	cityService := ht.NewServer(_http.HandlerFunc(cityServiceHandler))
	defer cityService.Close()

	populationData := map[int]int{
		1: 8_804_190,
		2: 3_898_747,
		3: 2_746_388,
	}

	populationServiceHandler := func(w _http.ResponseWriter, r *_http.Request) {
		r.ParseForm()
		if len(r.Form["id"]) > 0 {
			id, _ := strconv.Atoi(r.Form["id"][0])
			population := populationData[id]
			out := fmt.Sprintf(`{"id": %d, "population": %d}`, id, population)
			_, err := fmt.Fprintf(w, out)
			if err != nil {
				panic(err)
			}
		}
	}

	populationService := ht.NewServer(_http.HandlerFunc(populationServiceHandler))
	defer populationService.Close()

	weatherData := map[int]int{
		1: 46,
		2: 60,
		3: 54,
	}
	weatherServiceHandler := func(w _http.ResponseWriter, r *_http.Request) {
		r.ParseForm()
		if len(r.Form["id"]) > 0 {
			id, _ := strconv.Atoi(r.Form["id"][0])
			temperature := weatherData[id]
			out := fmt.Sprintf(`{"id": %d, "temperature": %d}`, id, temperature)
			_, err := fmt.Fprintf(w, out)
			if err != nil {
				panic(err)
			}
		}
	}

	weatherService := ht.NewServer(_http.HandlerFunc(weatherServiceHandler))
	defer weatherService.Close()
	//getCityReq := func [T any](cityName string) http.ReqRecipe[T] {
	//	return http.GetRecipe[T](cityService.URL, url.Values{"name": []string{cityName}})
	//}
	getCityReq := http.ReqRecipe[City](func(city *City) error {
		in := url.Values{}
		in.Set("name", city.Name)
		return http.Get(cityService.URL, in, &city)
	})

	getCityPopulationReq := http.ReqRecipe[City](func(city *City) error {
		in := url.Values{}
		in.Set("id", strconv.Itoa(city.Id))
		return http.Get(populationService.URL, url.Values{"id": []string{strconv.Itoa(city.Id)}}, &city)
	})

	getCityWeatherReq := http.ReqRecipe[City](func(city *City) error {
		in := url.Values{}
		in.Set("id", strconv.Itoa(city.Id))
		return http.Get(weatherService.URL, url.Values{"id": []string{strconv.Itoa(city.Id)}}, &city)
	})

	populateCity := func(city *City) error {
		return http.Pipe[City](city, getCityReq, getCityPopulationReq, getCityWeatherReq)
	}

	city := City{Name: "New York"}

	err := populateCity(&city)

	if err != nil {
		fmt.Println(err)
	}
	json.Object[City]{city}.Print()

	// Output:
	// {"id":1,"name":"New York","state":"New York","population":8804190,"temperature":46}
}

//func ExamplePipe() {
//
//	// setup mock web service
//	handler := _http.HandlerFunc(
//		func(w _http.ResponseWriter, r *_http.Request) {
//			switch r.URL.String() {
//			case "/greet":
//				fmt.Fprintf(w, `{"greeting":"hello"}`)
//			case "/howru":
//				fmt.Fprintf(w, `{"word":"hello","name":"Rob"}`)
//			case "/bye":
//				fmt.Fprintf(w, `{}`)
//			default:
//				w.WriteHeader(400)
//			}
//		})
//	svr := ht.NewServer(handler)
//	defer svr.Close()
//
//	greet := http.GET{svr.URL + "/greet", nil}
//	howru := http.POST{svr.URL + "/howru", url.Values{"name": []string{"Rob"}}}
//	bye := http.GET{svr.URL + "/bye", nil}
//
//	//encounter := []http.Req{greet, howru, bye}
//
//	type Data struct {
//		Name     string
//		Greeting string
//	}
//	data := Data{}
//
//	err := http.Pipe(&data, greet, howru, bye)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	json.Object[Data]{data}.Print()
//
//	// Output:
//	// {"Name":"Rob","Greeting":"hello"}
//}

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

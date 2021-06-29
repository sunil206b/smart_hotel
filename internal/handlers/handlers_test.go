package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name          string
	url           string
	method        string
	params        []postData
	expStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservations", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "POST", []postData{
		{key: "startDate", value: "2021-01-01"},
		{key: "endDate", value: "2021-01-10"},
	}, http.StatusOK},
	{"search-availability-json", "/search-availability-json", "POST", []postData{
		{key: "startDate", value: "2021-01-01"},
		{key: "endDate", value: "2021-01-10"},
	}, http.StatusOK},
	{"make-reservations", "/make-reservations", "POST", []postData{
		{key: "firstName", value: "John"},
		{key: "lastName", value: "Smith"},
		{key: "email", value: "smith@test.com"},
		{key: "phone", value: "767-432-4312"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			res, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if res.StatusCode != e.expStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expStatusCode, res.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, p := range e.params {
				values.Add(p.key, p.value)
			}
			res, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if res.StatusCode != e.expStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expStatusCode, res.StatusCode)
			}
		}
	}
}

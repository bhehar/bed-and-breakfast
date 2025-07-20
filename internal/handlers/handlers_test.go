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
	name               string
	url                string
	mehtod             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"make-reservation-post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "James"},
		{key: "last_name", value: "Harris"},
		{key: "email", value: "wise@fool.com"},
		{key: "phone", value: "661-383-3355"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	mux := getRoutes()
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	for _, tc := range theTests {
		if tc.mehtod == "GET" {
			resp, err := ts.Client().Get(ts.URL + tc.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("for %s, expect status %d got %d", tc.name, tc.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range tc.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+tc.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("for %s, expect status %d got %d", tc.name, tc.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

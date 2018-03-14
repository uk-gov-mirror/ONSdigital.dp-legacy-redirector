package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type urlTest struct {
	url      string
	code     int
	body     string
	location string
}

var tests = []urlTest{
	// Websites
	{"https://web.ons.gov.uk/", redir, "", landingPage},
	{"https://web.ons.gov.uk/a/b/c", redir, "", landingPage},
	{"https://web.ons.gov.uk/ons/apiservice/web/", redir, "", landingPage},
	// APIs
	{"https://neighbourhood.statistics.gov.uk/NDE2/a/b/c", 410, apiResponse, ""},
	{"https://web.ons.gov.uk/ons/apiservice/a/b/c", 410, apiResponse, ""},
	{"https://web.ons.gov.uk/ons/api/a/b/c", 410, apiResponse, ""},
	{"https://data.ons.gov.uk/ons/api/a/b/c", 410, apiResponse, ""},
	// Visualisations
	{"https://neighbourhood.statistics.gov.uk/HTMLDocs/a/b/c", redir, "", "https://www.ons.gov.uk/visualisations/nesscontent/a/b/c"},
	{"https://www.neighbourhood.statistics.gov.uk/HTMLDocs/a/b/c", redir, "", "https://www.ons.gov.uk/visualisations/nesscontent/a/b/c"},
	{"https://visual.ons.gov.uk/a/b/c", 410, "FIXME", ""},
}

func TestRedirects(t *testing.T) {
	router := getRouter()

	for _, test := range tests {
		Convey(fmt.Sprintf("%s", test.url), t, func() {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", test.url, nil)
			So(err, ShouldBeNil)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, test.code)
			So(w.Header().Get("Location"), ShouldEqual, test.location)
			So(w.Body.String(), ShouldEqual, test.body)
		})
	}
}

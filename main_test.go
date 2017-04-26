package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultHandler(t *testing.T) {
	router := getRouter()

	Convey("Default handler should redirect to landing page", t, func() {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "https://web.ons.gov.uk/a/b/c", nil)
		So(err, ShouldBeNil)

		router.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, redir)
		So(w.Header().Get("Location"), ShouldEqual, landingPage)
	})

	Convey("Data vis handler should redirect to /visualisations/nesscontent", t, func() {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "https://neighbourhood.statistics.gov.uk/HTMLDocs/a/b/c", nil)
		So(err, ShouldBeNil)

		router.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, redir)
		So(w.Header().Get("Location"), ShouldEqual, "https://www.ons.gov.uk/visualisations/nesscontent/a/b/c")
	})

	Convey("NeSS API handler should redirect to /help/localstatistics", t, func() {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "https://neighbourhood.statistics.gov.uk/NDE2/a/b/c", nil)
		So(err, ShouldBeNil)

		router.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
		So(w.Body.String(), ShouldEqual, `This service has been retired. Please visit https://www.ons.gov.uk/help/localstatistics for more information.`)
	})
}

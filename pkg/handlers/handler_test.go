package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/version"

	"github.com/takama/bit"
)

func TestRoot(t *testing.T) {
	h := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Root)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, fmt.Sprintf("%s %s", config.ServiceName, version.RELEASE))
}

func TestNotFound(t *testing.T) {
	h := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.NotFound)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusNotFound, "Method not found for /")
}

func testHandler(t *testing.T, handler http.HandlerFunc, code int, body string) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	trw := httptest.NewRecorder()
	handler.ServeHTTP(trw, req)

	if trw.Code != code {
		t.Error("Expected status code:", code, "got", trw.Code)
	}
	if trw.Body.String() != body {
		t.Error("Expected body", body, "got", trw.Body.String())
	}
}

func TestCollectCodes(t *testing.T) {
	h := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(func(c bit.Control) {
			c.Code(http.StatusBadGateway)
			c.Body(http.StatusText(http.StatusBadGateway))
		})(bit.NewControl(w, r))
	})
	testHandler(t, handler, http.StatusBadGateway, http.StatusText(http.StatusBadGateway))

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(func(c bit.Control) {
			c.Code(http.StatusNotFound)
			c.Body(http.StatusText(http.StatusNotFound))
		})(bit.NewControl(w, r))
	})
	testHandler(t, handler, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

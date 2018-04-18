package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/back-friend/pkg/version"

	"github.com/takama/bit"
)

var ErrTestError = errors.New("Test Error")

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

func testHandlerWithParams(t *testing.T, params map[string]string,
	handler *Handler, control func(bit.Control), code int, body string) {
	wrapHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		for idx, val := range params {
			ctrl.Params().Set(idx, val)
		}
		handler.Base(control)(ctrl)
	})
	testHandler(t, wrapHandler, code, body)
}

func TestRoot(t *testing.T) {
	conn := &db.Connection{
		Config:     config.New(),
		Controller: new(db.Mock),
	}
	h := New(conn)
	testHandlerWithParams(t,
		nil,
		h, h.Root,
		http.StatusOK, fmt.Sprintf("%s %s", config.ServiceName, version.RELEASE))
}

func TestNotFound(t *testing.T) {
	conn := &db.Connection{
		Config:     config.New(),
		Controller: new(db.Mock),
	}
	h := New(conn)
	testHandlerWithParams(t,
		nil,
		h, h.NotFound,
		http.StatusNotFound, "Method not found for /")
}

func TestCollectCodes(t *testing.T) {
	conn := &db.Connection{
		Config:     config.New(),
		Controller: new(db.Mock),
	}
	h := New(conn)
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

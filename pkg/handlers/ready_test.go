package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/bit"
)

func TestReady(t *testing.T) {
	conn := &db.Connection{
		Driver: db.Stub{},
	}
	h := New(conn)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Ready)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, http.StatusText(http.StatusOK))
}

func TestNotReady(t *testing.T) {
	conn := &db.Connection{
		Driver: db.Mock{
			OnReady: func() bool { return false },
		},
	}
	h := New(conn)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Ready)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusServiceUnavailable,
		http.StatusText(http.StatusServiceUnavailable))
}

package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/bit"
)

func TestReady(t *testing.T) {
	mock := &db.Mock{
		OnReady: func() bool { return true },
	}
	conn := &db.Connection{
		Config:     config.New(),
		Controller: mock,
		Store:      mock,
	}
	h := New(conn)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Ready)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, http.StatusText(http.StatusOK))
}

func TestNotReady(t *testing.T) {
	mock := &db.Mock{
		OnReady: func() bool { return false },
	}
	conn := &db.Connection{
		Config:     config.New(),
		Controller: mock,
		Store:      mock,
	}
	h := New(conn)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Ready)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusServiceUnavailable,
		http.StatusText(http.StatusServiceUnavailable))
}

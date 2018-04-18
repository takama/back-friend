package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
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
	testHandlerWithParams(t,
		nil,
		h, h.Ready,
		http.StatusOK, http.StatusText(http.StatusOK))
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
	testHandlerWithParams(t,
		nil,
		h, h.Ready,
		http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
}

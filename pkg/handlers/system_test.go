package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
)

var ErrResetFalse = errors.New("Reset false")

func TestReset(t *testing.T) {
	mock := &db.Mock{
		OnReset: func() error { return nil },
	}
	conn := &db.Connection{
		Config:     config.New(),
		Controller: mock,
		Store:      mock,
	}
	h := New(conn)
	testHandlerWithParams(t,
		nil,
		h, h.Reset,
		http.StatusOK, http.StatusText(http.StatusOK))
}

func TestResetFalse(t *testing.T) {
	mock := &db.Mock{
		OnReset: func() error { return ErrResetFalse },
	}
	conn := &db.Connection{
		Config:     config.New(),
		Controller: mock,
		Store:      mock,
	}
	h := New(conn)
	testHandlerWithParams(t,
		nil,
		h, h.Reset,
		http.StatusInternalServerError, ErrResetFalse.Error())
}

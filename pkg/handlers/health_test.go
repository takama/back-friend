package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
)

func TestHealth(t *testing.T) {
	conn := &db.Connection{
		Config:     config.New(),
		Controller: new(db.Mock),
	}
	h := New(conn)
	testHandlerWithParams(t,
		nil,
		h, h.Health,
		http.StatusOK, http.StatusText(http.StatusOK))
}

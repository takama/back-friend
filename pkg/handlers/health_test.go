package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/bit"
)

func TestHealth(t *testing.T) {
	conn := &db.Connection{
		Config:     config.New(),
		Controller: new(db.Mock),
	}
	h := New(conn)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Health)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, http.StatusText(http.StatusOK))
}

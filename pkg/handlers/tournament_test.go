package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/backer/datastore"
	"github.com/takama/bit"
)

func TestTournamentDetails(t *testing.T) {
	stub := new(datastore.Stub)
	stub.Reset()
	conn := &db.Connection{
		Config:     config.New(),
		Controller: stub,
		Store:      stub,
	}
	h := New(conn)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		ctrl.Params().Set(":id", "1")
		h.Base(h.TournamentDetails)(ctrl)
	})

	testHandler(t, handler, http.StatusNotFound, datastore.ErrRecordNotFound.Error())
}

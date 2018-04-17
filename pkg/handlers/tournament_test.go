package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/backer/datastore"
	"github.com/takama/bit"
)

var ErrTestError = errors.New("Test Error")

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

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		ctrl.Params().Set(":id", "0")
		h.Base(h.TournamentDetails)(ctrl)
	})

	testHandler(t, handler, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	stub.NewTournament(1, nil)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		ctrl.Params().Set(":id", "1")
		h.Base(h.TournamentDetails)(ctrl)
	})

	testHandler(t, handler, http.StatusOK, `{"id":1,"deposit":0,"is_finished":false,"bidders":[]}`)

	stub.ErrFind = append(stub.ErrFind, ErrTestError)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		ctrl.Params().Set(":id", "1")
		h.Base(h.TournamentDetails)(ctrl)
	})

	testHandler(t, handler, http.StatusInternalServerError, ErrTestError.Error())
}

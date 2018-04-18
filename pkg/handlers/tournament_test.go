package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/backer/datastore"
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

	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentDetails,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "0"},
		h, h.TournamentDetails,
		http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	stub.NewTournament(1, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentDetails,
		http.StatusOK, `{"id":1,"deposit":0,"is_finished":false,"bidders":[]}`)

	stub.ErrFind = append(stub.ErrFind, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentDetails,
		http.StatusInternalServerError, ErrTestError.Error())
}

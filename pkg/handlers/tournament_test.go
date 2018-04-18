package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/backer/datastore"
	"github.com/takama/backer/model"
	"github.com/takama/backer/tournament"
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

func TestTournamentAnnounce(t *testing.T) {
	stub := new(datastore.Stub)
	stub.Reset()
	conn := &db.Connection{
		Config:     config.New(),
		Controller: stub,
		Store:      stub,
	}
	h := New(conn)

	testHandlerWithParams(t,
		map[string]string{":id": "0"},
		h, h.TournamentAnnounce,
		http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	stub.ErrNew = append(stub.ErrNew, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentAnnounce,
		http.StatusInternalServerError, ErrTestError.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentAnnounce,
		http.StatusBadRequest, couldNotRecognizeRequestData)

	stub.ErrFind = append(stub.ErrFind, datastore.ErrRecordNotFound, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"deposit": 300}`},
		h, h.TournamentAnnounce,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.ErrTx = append(stub.ErrFind, ErrTestError, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"deposit": 300}`},
		h, h.TournamentAnnounce,
		http.StatusInternalServerError, ErrTestError.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"deposit": 300}`},
		h, h.TournamentAnnounce,
		http.StatusOK, `{"id":1,"deposit":300,"is_finished":false,"bidders":[]}`)

	tx, _ := stub.Transaction()
	stub.SaveTournament(&model.Tournament{ID: 1, IsFinished: true}, tx)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"deposit": 500}`},
		h, h.TournamentAnnounce,
		http.StatusNotAcceptable, tournament.ErrAllreadyFinished.Error())
}

func TestTournamentJoin(t *testing.T) {
	stub := new(datastore.Stub)
	stub.Reset()
	conn := &db.Connection{
		Config:     config.New(),
		Controller: stub,
		Store:      stub,
	}
	h := New(conn)

	testHandlerWithParams(t,
		map[string]string{":id": "0"},
		h, h.TournamentJoin,
		http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentJoin,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.NewTournament(1, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentJoin,
		http.StatusBadRequest, couldNotRecognizeRequestData)

	stub.ErrFind = append(stub.ErrFind, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentJoin,
		http.StatusInternalServerError, ErrTestError.Error())
}

func TestTournamentResult(t *testing.T) {
	stub := new(datastore.Stub)
	stub.Reset()
	conn := &db.Connection{
		Config:     config.New(),
		Controller: stub,
		Store:      stub,
	}
	h := New(conn)

	testHandlerWithParams(t,
		map[string]string{":id": "0"},
		h, h.TournamentResult,
		http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentResult,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.NewTournament(1, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentResult,
		http.StatusBadRequest, couldNotRecognizeRequestData)

	stub.ErrFind = append(stub.ErrFind, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentResult,
		http.StatusInternalServerError, ErrTestError.Error())
}

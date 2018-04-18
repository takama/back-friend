package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/backer/datastore"
	"github.com/takama/backer/model"
	"github.com/takama/backer/player"
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

	stub.ErrTx = append(stub.ErrTx, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "1"},
		h, h.TournamentAnnounce,
		http.StatusInternalServerError, ErrTestError.Error())

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

	stub.SaveTournament(&model.Tournament{ID: 1, IsFinished: true}, nil)
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

	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1", "backers": ["b1", "b2"]}`},
		h, h.TournamentJoin,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.ErrTx = append(stub.ErrTx, ErrTestError, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1", "backers": ["b1", "b2"]}`},
		h, h.TournamentJoin,
		http.StatusInternalServerError, ErrTestError.Error())

	stub.NewPlayer("p1", nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1", "backers": ["b1", "b2"]}`},
		h, h.TournamentJoin,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.ErrTx = append(stub.ErrTx, ErrTestError, nil, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1", "backers": ["b1", "b2"]}`},
		h, h.TournamentJoin,
		http.StatusInternalServerError, ErrTestError.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"deposit": 300}`},
		h, h.TournamentAnnounce,
		http.StatusOK, `{"id":1,"deposit":300,"is_finished":false,"bidders":[]}`)

	stub.NewPlayer("b1", nil)
	stub.NewPlayer("b2", nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1", "backers": ["b1", "b2"]}`},
		h, h.TournamentJoin,
		http.StatusNotAcceptable, player.ErrInsufficientPoints.Error())

	stub.ErrFind = append(stub.ErrFind, datastore.ErrRecordNotFound, nil, nil, nil, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1", "backers": ["b1", "b2"]}`},
		h, h.TournamentJoin,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.ErrTx = append(stub.ErrTx, ErrTestError, nil, nil, nil, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1", "backers": ["b1", "b2"]}`},
		h, h.TournamentJoin,
		http.StatusInternalServerError, ErrTestError.Error())

	stub.SavePlayer(&model.Player{ID: "p1", Balance: 100}, nil)
	stub.SavePlayer(&model.Player{ID: "b1", Balance: 100}, nil)
	stub.SavePlayer(&model.Player{ID: "b2", Balance: 100}, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1", "backers": ["b1", "b2"]}`},
		h, h.TournamentJoin,
		http.StatusOK, `{"id":1,"deposit":300,"is_finished":false,"bidders":[{"id":"p1","winner":false,"prize":0,"backers":["b1","b2"]}]}`)
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

	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"winners": [{"player": "p1", "prize": 300}]}`},
		h, h.TournamentResult,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"deposit": 300}`},
		h, h.TournamentAnnounce,
		http.StatusOK, `{"id":1,"deposit":300,"is_finished":false,"bidders":[]}`)

	stub.SavePlayer(&model.Player{ID: "p1", Balance: 300}, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"player": "p1"}`},
		h, h.TournamentJoin,
		http.StatusOK, `{"id":1,"deposit":300,"is_finished":false,"bidders":[{"id":"p1","winner":false,"prize":0,"backers":[]}]}`)

	stub.ErrFind = append(stub.ErrFind, ErrTestError, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"winners": [{"player": "p1", "prize": 300}]}`},
		h, h.TournamentResult,
		http.StatusInternalServerError, ErrTestError.Error())

	stub.ErrFind = append(stub.ErrFind, datastore.ErrRecordNotFound, nil, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"winners": [{"player": "p1", "prize": 300}]}`},
		h, h.TournamentResult,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.ErrTx = append(stub.ErrTx, ErrTestError, nil, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"winners": [{"player": "p1", "prize": 300}]}`},
		h, h.TournamentResult,
		http.StatusInternalServerError, ErrTestError.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"winners": [{"player": "p1", "prize": 300}]}`},
		h, h.TournamentResult,
		http.StatusOK, `{"id":1,"deposit":300,"is_finished":true,"bidders":[{"id":"p1","winner":true,"prize":300,"backers":[]}]}`)

	testHandlerWithParams(t,
		map[string]string{":id": "1", "data": `{"winners": [{"player": "p1", "prize": 300}]}`},
		h, h.TournamentResult,
		http.StatusNotAcceptable, tournament.ErrAllreadyFinished.Error())
}

package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"

	"github.com/takama/backer/datastore"
	"github.com/takama/backer/player"
)

func TestPlayerDetails(t *testing.T) {
	stub := new(datastore.Stub)
	stub.Reset()
	conn := &db.Connection{
		Config:     config.New(),
		Controller: stub,
		Store:      stub,
	}
	h := New(conn)

	testHandlerWithParams(t,
		map[string]string{":id": ""},
		h, h.PlayerDetails,
		http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	testHandlerWithParams(t,
		map[string]string{":id": "not-existing"},
		h, h.PlayerDetails,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.NewPlayer("p1", nil)
	testHandlerWithParams(t,
		map[string]string{":id": "p1"},
		h, h.PlayerDetails,
		http.StatusOK, `{"id":"p1","balance":0}`)

	stub.ErrFind = append(stub.ErrFind, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "p1"},
		h, h.PlayerDetails,
		http.StatusInternalServerError, ErrTestError.Error())
}

func TestPlayerFund(t *testing.T) {
	stub := new(datastore.Stub)
	stub.Reset()
	conn := &db.Connection{
		Config:     config.New(),
		Controller: stub,
		Store:      stub,
	}
	h := New(conn)

	testHandlerWithParams(t,
		map[string]string{":id": ""},
		h, h.PlayerFund,
		http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	stub.ErrNew = append(stub.ErrNew, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "p1"},
		h, h.PlayerFund,
		http.StatusInternalServerError, ErrTestError.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "p1"},
		h, h.PlayerFund,
		http.StatusBadRequest, couldNotRecognizeRequestData)

	stub.ErrFind = append(stub.ErrFind, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "p1"},
		h, h.PlayerFund,
		http.StatusInternalServerError, ErrTestError.Error())

	stub.ErrFind = append(stub.ErrFind, ErrTestError, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "p1", "data": `{"points": 300}`},
		h, h.PlayerFund,
		http.StatusInternalServerError, ErrTestError.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "p1", "data": `{"points": 300}`},
		h, h.PlayerFund,
		http.StatusOK, `{"id":"p1","balance":300}`)
}

func TestPlayerTake(t *testing.T) {
	stub := new(datastore.Stub)
	stub.Reset()
	conn := &db.Connection{
		Config:     config.New(),
		Controller: stub,
		Store:      stub,
	}
	h := New(conn)

	testHandlerWithParams(t,
		map[string]string{":id": ""},
		h, h.PlayerTake,
		http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	testHandlerWithParams(t,
		map[string]string{":id": "p1"},
		h, h.PlayerTake,
		http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	stub.NewPlayer("p1", nil)
	testHandlerWithParams(t,
		map[string]string{":id": "p1"},
		h, h.PlayerTake,
		http.StatusBadRequest, couldNotRecognizeRequestData)

	stub.ErrFind = append(stub.ErrFind, ErrTestError)
	testHandlerWithParams(t,
		map[string]string{":id": "p1"},
		h, h.PlayerTake,
		http.StatusInternalServerError, ErrTestError.Error())

	stub.ErrFind = append(stub.ErrFind, ErrTestError, nil)
	testHandlerWithParams(t,
		map[string]string{":id": "p1", "data": `{"points": 300}`},
		h, h.PlayerTake,
		http.StatusInternalServerError, ErrTestError.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "p1", "data": `{"points": 300}`},
		h, h.PlayerTake,
		http.StatusNotAcceptable, player.ErrInsufficientPoints.Error())

	testHandlerWithParams(t,
		map[string]string{":id": "p1", "data": `{"points": 300}`},
		h, h.PlayerFund,
		http.StatusOK, `{"id":"p1","balance":300}`)

	testHandlerWithParams(t,
		map[string]string{":id": "p1", "data": `{"points": 300}`},
		h, h.PlayerTake,
		http.StatusOK, `{"id":"p1","balance":0}`)
}

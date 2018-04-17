package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"

	"github.com/takama/backer/datastore"
	"github.com/takama/bit"
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
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		ctrl.Params().Set(":id", "not-existing")
		h.Base(h.PlayerDetails)(ctrl)
	})

	testHandler(t, handler, http.StatusNotFound, datastore.ErrRecordNotFound.Error())

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		ctrl.Params().Set(":id", "")
		h.Base(h.PlayerDetails)(ctrl)
	})

	testHandler(t, handler, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

	stub.NewPlayer("p1", nil)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		ctrl.Params().Set(":id", "p1")
		h.Base(h.PlayerDetails)(ctrl)
	})

	testHandler(t, handler, http.StatusOK, `{"id":"p1","balance":0}`)

	stub.ErrFind = append(stub.ErrFind, ErrTestError)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := bit.NewControl(w, r)
		ctrl.Params().Set(":id", "p1")
		h.Base(h.PlayerDetails)(ctrl)
	})

	testHandler(t, handler, http.StatusInternalServerError, ErrTestError.Error())
}

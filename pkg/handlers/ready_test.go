package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/back-friend/pkg/logger"
	"github.com/takama/back-friend/pkg/logger/stdlog"
	"github.com/takama/bit"
)

func TestReady(t *testing.T) {
	conn, _, _ := db.New(config.New(), stdlog.New(new(logger.Config)))
	h := New(conn)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Ready)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, http.StatusText(http.StatusOK))
}

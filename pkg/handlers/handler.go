package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/back-friend/pkg/version"

	"github.com/takama/bit"
)

// Handler defines common part for all handlers
type Handler struct {
	maintenance bool
	db          *db.Connection
	stats       *stats
}

type stats struct {
	requests  *Requests
	startTime time.Time
}

// New returns new instance of the Handler
func New(conn *db.Connection) *Handler {
	return &Handler{
		db: conn,
		stats: &stats{
			requests:  new(Requests),
			startTime: time.Now(),
		},
	}
}

// Base handler implements middleware logic
func (h *Handler) Base(handle func(bit.Control)) func(bit.Control) {
	return func(c bit.Control) {
		handle(c)
		h.collectCodes(c)
	}
}

// Root handler shows version
func (h *Handler) Root(c bit.Control) {
	c.Code(http.StatusOK)
	c.Body(fmt.Sprintf("%s %s", config.ServiceName, version.RELEASE))
}

// NotFound responds for undefined methods
func (h *Handler) NotFound(c bit.Control) {
	c.Code(http.StatusNotFound)
	c.Body("Method not found for " + c.Request().URL.Path)
}

func (h *Handler) collectCodes(c bit.Control) {
	if c.GetCode() >= 500 {
		h.stats.requests.C5xx++
	} else {
		if c.GetCode() >= 400 {
			h.stats.requests.C4xx++
		} else {
			if c.GetCode() >= 200 && c.GetCode() < 300 {
				h.stats.requests.C2xx++
			}
		}
	}
}

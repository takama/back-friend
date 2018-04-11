package handlers

import (
	"net/http"

	"github.com/takama/bit"
)

// Ready returns "OK" if service is ready to serve traffic
func (h *Handler) Ready(c bit.Control) {
	// TODO: possible use cases:
	// load data from a database, a message broker, any external services, etc

	c.Code(http.StatusOK)
	c.Body(http.StatusText(http.StatusOK))
}

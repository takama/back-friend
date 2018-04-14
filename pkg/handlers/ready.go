package handlers

import (
	"net/http"

	"github.com/takama/bit"
)

// Ready returns "OK" if service is ready to serve traffic
func (h *Handler) Ready(c bit.Control) {
	if !h.db.Ready() {
		c.Code(http.StatusServiceUnavailable)
		c.Body(http.StatusText(http.StatusServiceUnavailable))
		return
	}
	c.Code(http.StatusOK)
	c.Body(http.StatusText(http.StatusOK))
}

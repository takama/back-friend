package handlers

import (
	"net/http"

	"github.com/takama/bit"
)

// Health returns "OK" if service is alive
func (h *Handler) Health(c bit.Control) {
	c.Code(http.StatusOK)
	c.Body(http.StatusText(http.StatusOK))
}

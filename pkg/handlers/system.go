package handlers

import (
	"net/http"

	"github.com/takama/bit"
)

// Reset makes the DB initialization
func (h *Handler) Reset(c bit.Control) {
	err := h.db.Reset()
	if err != nil {
		serviceError(err, c)
		return
	}
	c.Code(http.StatusOK)
	c.Body(http.StatusText(http.StatusOK))
}

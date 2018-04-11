package handlers

import (
	"fmt"
	"net/http"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/version"

	"github.com/takama/bit"
)

// Handler defines common part for all handlers
type Handler struct {
}

// Base handler implements middleware logic
func (h *Handler) Base(handle func(bit.Control)) func(bit.Control) {
	return func(c bit.Control) {
		handle(c)
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

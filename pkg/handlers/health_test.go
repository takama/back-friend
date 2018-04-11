package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/bit"
)

func TestHealth(t *testing.T) {
	h := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Health)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, http.StatusText(http.StatusOK))
}

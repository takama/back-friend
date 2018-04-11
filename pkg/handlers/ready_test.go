package handlers

import (
	"net/http"
	"testing"

	"github.com/takama/bit"
)

func TestReady(t *testing.T) {
	h := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Ready)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, http.StatusText(http.StatusOK))
}

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takama/back-friend/pkg/config"
	"github.com/takama/back-friend/pkg/db"
	"github.com/takama/back-friend/pkg/version"

	"github.com/takama/bit"
)

func TestInfo(t *testing.T) {
	conn, _ := db.New(config.New())
	h := New(conn)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Info)(bit.NewControl(w, r))
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	trw := httptest.NewRecorder()
	handler.ServeHTTP(trw, req)

	if trw.Code != http.StatusOK {
		t.Error("Expected status:", http.StatusOK, "got", trw.Code)
	}

	var s Status
	err = json.Unmarshal(trw.Body.Bytes(), &s)
	if err != nil {
		t.Fatal(err)
	}

	if s.Version != version.RELEASE {
		t.Error("Expected version:", version.RELEASE, "got", s.Version)
	}
	if s.Commit != version.COMMIT {
		t.Error("Expected commit hash:", version.COMMIT, "got", s.Commit)
	}
	if s.Repo != version.REPO {
		t.Error("Expected repository:", version.REPO, "got", s.Repo)
	}
}

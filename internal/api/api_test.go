package api

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"independentweeklylog/internal/query"
	"independentweeklylog/internal/repository"
	"independentweeklylog/internal/review"
	"independentweeklylog/internal/store"
	"independentweeklylog/internal/workflow20"
)

func TestHealthEndpointReturnsStatus(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "journal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	server := NewServer(workflow20.New(repo), review.New(repo), query.New(repo))
	req := httptest.NewRequest("GET", "/health", nil)
	res := httptest.NewRecorder()
	server.Handler().ServeHTTP(res, req)
	if res.Code != 200 || !strings.Contains(res.Body.String(), "ok") {
		t.Fatalf("unexpected response: %d %s", res.Code, res.Body.String())
	}
}

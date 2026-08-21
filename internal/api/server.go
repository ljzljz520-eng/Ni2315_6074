package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/query"
	"independentweeklylog/internal/repository"
	"independentweeklylog/internal/review"
	"independentweeklylog/internal/workflow20"
)

type Server struct {
	workflow *workflow20.Orchestrator
	review   *review.Service
	query    *query.Service
}

func NewServer(workflow *workflow20.Orchestrator, reviewer *review.Service, finder *query.Service) *Server {
	return &Server{workflow: workflow, review: reviewer, query: finder}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/entries", s.entries)
	mux.HandleFunc("/summary", s.summary)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) entries(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		week, _ := strconv.Atoi(req.URL.Query().Get("week"))
		items, err := s.query.Search(repository.EntryFilter{Week: week, Author: req.URL.Query().Get("author"), Tag: req.URL.Query().Get("tag"), Text: req.URL.Query().Get("q")})
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	if req.Method != http.MethodPost {
		writeErrorStatus(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var entry domain.JournalEntry
	if err := json.NewDecoder(req.Body).Decode(&entry); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, "invalid json")
		return
	}
	result, err := s.workflow.CaptureDraft(entry, nil)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result.Entry)
}

func (s *Server) summary(w http.ResponseWriter, req *http.Request) {
	week, _ := strconv.Atoi(req.URL.Query().Get("week"))
	summary, err := s.query.Summary(week)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	writeErrorStatus(w, http.StatusBadRequest, err.Error())
}

func writeErrorStatus(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": strings.TrimSpace(message)})
}

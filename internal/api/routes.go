package api

import (
	"net/http"

	"independentweeklylog/internal/archive"
	"independentweeklylog/internal/repository"
)

func ArchiveHandler(repo *repository.Repository, archiver *archive.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			writeErrorStatus(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		entryID := req.URL.Query().Get("entry_id")
		record, err := archiver.Archive(entryID, "weekly close")
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, record)
	}
}

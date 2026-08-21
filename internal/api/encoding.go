package api

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Data  any            `json:"data,omitempty"`
	Error string         `json:"error,omitempty"`
	Meta  map[string]any `json:"meta,omitempty"`
}

func writeEnvelope(w http.ResponseWriter, status int, data any, meta map[string]any) {
	writeJSON(w, status, envelope{Data: data, Meta: meta})
}

func decodeBody(req *http.Request, target any) error {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func acceptsJSON(req *http.Request) bool {
	return req.Header.Get("Accept") == "" || req.Header.Get("Accept") == "application/json"
}

func requestActor(req *http.Request) string { return req.Header.Get("X-Actor") }

func pageBounds(req *http.Request, total int) (int, int) {
	page, size := 1, 20
	if raw := req.URL.Query().Get("page"); raw != "" {
		if parsed := atoi(raw); parsed > 0 {
			page = parsed
		}
	}
	if raw := req.URL.Query().Get("size"); raw != "" {
		if parsed := atoi(raw); parsed > 0 && parsed <= 100 {
			size = parsed
		}
	}
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}
	return start, end
}

func atoi(value string) int {
	number := 0
	for _, char := range value {
		if char < '0' || char > '9' {
			return 0
		}
		number = number*10 + int(char-'0')
	}
	return number
}

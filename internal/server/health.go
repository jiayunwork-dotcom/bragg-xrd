package server

import (
	"net/http"
	"runtime"
	"time"
)

var startedAt = time.Now()

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok", "uptime": int64(time.Since(startedAt).Seconds()),
		"goos": runtime.GOOS, "goarch": runtime.GOARCH, "version": "1.0.0",
	})
}

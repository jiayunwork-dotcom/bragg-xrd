package server

import (
	"encoding/json"
	"net/http"
)

func Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/bragg", braggHandler)
	mux.HandleFunc("/api/powder", powderHandler)
	mux.HandleFunc("/api/version", versionHandler)
	return mux
}

func readJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	if r.Body == nil {
		writeError(w, http.StatusBadRequest, "missing request body")
		return false
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return false
	}
	return true
}

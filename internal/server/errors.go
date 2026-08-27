package server

import "net/http"

type errorBody struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorBody{Error: message})
}

func badRequest(w http.ResponseWriter, message string) {
	writeError(w, http.StatusBadRequest, message)
}

func methodNotAllowed(w http.ResponseWriter, method string) {
	writeError(w, http.StatusMethodNotAllowed, "method must be "+method)
}

func internalError(w http.ResponseWriter, message string) {
	writeError(w, http.StatusInternalServerError, message)
}

func notFound(w http.ResponseWriter, message string) {
	writeError(w, http.StatusNotFound, message)
}

func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				internalError(w, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

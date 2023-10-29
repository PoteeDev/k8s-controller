package handlers

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, err string, code int) {
	JSONResponse(w, map[string]string{
		"error": err,
	}, code)
}

func JSONResponse(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

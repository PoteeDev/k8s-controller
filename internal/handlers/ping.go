package handlers

import (
	"encoding/json"
	"net/http"
)

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"answer": "ping",
	})
}

package handlers

import (
	"net/http"

	client "github.com/PoteeDev/k8s-controller/internal/k8s_client"
)

func (s *Server) DestroyStand(w http.ResponseWriter, r *http.Request) {
	namespace, err := s.GenerateNamespace(r)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	chart, err := client.InitChart("", namespace)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = chart.Destroy(namespace)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]interface{}{
		"answer": map[string]string{
			"action":    "destroyed",
			"namespace": namespace,
		},
	}, http.StatusAccepted)
}

package handlers

import (
	"encoding/json"
	"net/http"

	client "github.com/PoteeDev/k8s-controller/internal/k8s_client"
)

type Stand struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (s *Server) DeployStand(w http.ResponseWriter, r *http.Request) {
	var stand Stand
	err := json.NewDecoder(r.Body).Decode(&stand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chart := client.InitChart(stand.Name, "", "")
	// validate incoming chart
	if err = chart.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// update repo and download chart if not exists
	if err = chart.AddRepo(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// deploy chart
	if err = chart.Deploy(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	stand.Namespace = chart.Namespace
	json.NewEncoder(w).Encode(&stand)
}

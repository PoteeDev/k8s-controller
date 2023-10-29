package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	client "github.com/PoteeDev/k8s-controller/internal/k8s_client"
	"github.com/PoteeDev/k8s-controller/internal/utils"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type Stand struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	ExtraValues map[string]string `json:"extra_values"`
}

func GenerateNamespace(response *oidc.IntrospectionResponse) (string, error) {
	log.Println(response)
	namespace := fmt.Sprintf("stand-%s", response.Subject)
	return namespace, nil
}

func (s *Server) DeployStand(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.ExtractToken(r)
	log.Println(token)
	resp, ierr := rs.Introspect[*oidc.IntrospectionResponse](r.Context(), s.ResourceServer, token)
	if ierr != nil {
		JSONError(w, ierr.Error(), http.StatusForbidden)
		return
	}
	namespace, err := GenerateNamespace(resp)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var stand Stand
	err = json.NewDecoder(r.Body).Decode(&stand)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	chart, err := client.InitChart(stand.Name, namespace)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validate incoming chart
	if err = chart.Validate(); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// update repo and download chart if not exists
	if err = chart.AddRepo(); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(stand.ExtraValues) == 0 {
		stand.ExtraValues = make(map[string]string)
	}
	stand.ExtraValues["flag"] = utils.GenerateFlag(30)

	// deploy chart
	if err = chart.Deploy(stand.ExtraValues); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	stand.Namespace = chart.Spec.Namespace
	JSONResponse(w, &stand, http.StatusAccepted)
}

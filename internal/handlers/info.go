package handlers

import (
	"log"
	"net/http"

	client "github.com/PoteeDev/k8s-controller/internal/k8s_client"
	"github.com/PoteeDev/k8s-controller/internal/utils"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func (s *Server) InfoHandler(w http.ResponseWriter, r *http.Request) {
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
	info, err := client.GetNamespaceInfo(namespace)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	JSONResponse(w, info, http.StatusOK)
}

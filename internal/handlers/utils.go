package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PoteeDev/k8s-controller/internal/utils"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"
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

func (s *Server) GenerateNamespace(r *http.Request) (string, error) {
	token, err := utils.ExtractToken(r)
	if err != nil {
		return "", err
	}
	log.Println(token)
	resp, err := rs.Introspect[*oidc.IntrospectionResponse](r.Context(), s.ResourceServer, token)
	if err != nil {
		return "", err
	}
	namespace := fmt.Sprintf("stand-%s", resp.Subject)
	return namespace, nil

}

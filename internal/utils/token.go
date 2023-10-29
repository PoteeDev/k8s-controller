package utils

import (
	"net/http"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func ExtractToken(r *http.Request) (string, error) {
	auth := r.Header.Get("authorization")
	return strings.TrimPrefix(auth, oidc.PrefixBearer), nil
}

package handlers

import (
	"github.com/zitadel/oidc/v3/pkg/client/rs"
)

type Server struct {
	ResourceServer rs.ResourceServer
}

func InitServer(ts rs.ResourceServer) *Server {
	s := &Server{ResourceServer: ts}
	return s
}

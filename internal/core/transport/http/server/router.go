package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/DenisAstrakhan/api-server/internal/core/transport/http/middleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	middleware []core_http_middleware.Middlware
}

func NewAPIVersionRouter(apiVersion ApiVersion, middleware ...core_http_middleware.Middlware) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.WithMiddlware())
	}
}
func (r *APIVersionRouter) WithMiddlware() http.Handler {
	return core_http_middleware.ChainMiddlare(r, r.middleware...)
}

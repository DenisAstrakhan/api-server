package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/DenisAstrakhan/api-server/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string                           //POST, GET, DELETe, PATCH, PUL ....
	Path       string                           // путь (/tasks или /users ...)
	Handler    http.HandlerFunc                 // который надо вызывать при заданно методе и пути
	Middleware []core_http_middleware.Middlware //набор middlware
}

func (r *Route) WithMiddlware() http.Handler {
	return core_http_middleware.ChainMiddlare(r.Handler, r.Middleware...)
}

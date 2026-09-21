package core_http_server

import "net/http"

type Route struct {
	Method  string           //POST, GET, DELETe, PATCH, PUL ....
	Path    string           // путь (/tasks или /users ...)
	Handler http.HandlerFunc // который надо вызывать при заданно методе и пути
}

func NewRoute(method string, path string, handler http.HandlerFunc) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}

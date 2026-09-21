package core_http_middleware

import "net/http"

type Middlware func(http.Handler) http.Handler

func ChainMiddlare(h http.Handler, m ...Middlware) http.Handler {
	if len(m) == 0 {
		return h
	}
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

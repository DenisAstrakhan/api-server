package core_http_middleware

import (
	"fmt"
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
)

func Dummy(s string) Middlware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			log.Debug(fmt.Sprintf("---> before: %s", s))

			next.ServeHTTP(w, r) //вызываем функцию которую оборачивает наша middleware
			log.Debug(fmt.Sprintf("<--- after: %s", s))
		})
	}
}

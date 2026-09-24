package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middlware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString() //если ID запроса не проставлен генерируем его
			}
			r.Header.Set(requestIDHeader, requestID)   // записываем ID запроса в хедер запроса
			w.Header().Set(requestIDHeader, requestID) // записываем ID запроса в хедер ответа
			next.ServeHTTP(w, r)                       //вызываем функцию которую оборачивает наша middleware
		})
	}
}

func Loger(log *core_logger.Logge) Middlware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)
			ctx := core_logger.ToContext(r.Context(), l) //context.WithValue(r.Context(), core_logger.LoggerContextKey, l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

func Trace() Middlware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w) //создаём свой ResponseWriter
			before := time.Now()
			log.Debug(
				"incoming HTTP request",
				zap.String("http method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
			)

		})
	}
}

func Panic() Middlware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

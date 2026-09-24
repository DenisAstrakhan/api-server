package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_middleware "github.com/DenisAstrakhan/api-server/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     config
	log        *core_logger.Logge
	middleware []core_http_middleware.Middlware
}

func NewHTTPServer(config config, log *core_logger.Logge, middleware ...core_http_middleware.Middlware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.WithMiddlware()))
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	//навешиваем middleware на мультеплексер
	mux := core_http_middleware.ChainMiddlare(s.mux, s.middleware...)
	//создаём сервер
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}
	// потдержка плавного завершения
	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		s.log.Warn("start HTTP server", zap.String("addr", s.config.Addr))
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		return fmt.Errorf("listen and serve HTTP: %w", err)
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server ...")

		shutdownCtx, cansel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cansel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		s.log.Warn("HTTP server stopped")
	}
	return nil
}

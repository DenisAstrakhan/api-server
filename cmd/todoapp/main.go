package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_pgx_pool "github.com/DenisAstrakhan/api-server/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/DenisAstrakhan/api-server/internal/core/transport/http/middleware"
	core_http_server "github.com/DenisAstrakhan/api-server/internal/core/transport/http/server"
	user_postgres_repository "github.com/DenisAstrakhan/api-server/internal/feature/users/repository/postgres"
	users_service "github.com/DenisAstrakhan/api-server/internal/feature/users/service"
	users_transport_http "github.com/DenisAstrakhan/api-server/internal/feature/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	// создаём контекст завязанный на системные сигналы
	ctx, cansel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cansel()

	logerConfig := core_logger.NewConfigMust()
	log, err := core_logger.NewLogger(logerConfig)
	if err != nil {
		fmt.Printf("логер не собрался: %v", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Debug("initiazling conection pool")
	//pool, err := core_postgres_pool.NewConnectionPool(core_postgres_pool.NewConfigMast(), ctx)
	pool, err := core_pgx_pool.NewPool(core_pgx_pool.NewConfigMast(), ctx)
	if err != nil {
		log.Fatal("failed to init connection pool", zap.Error(err))
	}
	defer pool.Close()

	log.Debug("initiazling feature", zap.String("feature", "users"))
	usersRepository := user_postgres_repository.NewUsersRepository(pool)
	userService := users_service.MewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(&userService)

	log.Debug("initiazling HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		log,
		core_http_middleware.RequestID(),
		core_http_middleware.Loger(log),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)

	//apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion2, core_http_middleware.Dummy("API v2 middleware"))
	//apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1 /*apiVersionRouterV2*/)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error: %w", zap.Error(err))
	}
}

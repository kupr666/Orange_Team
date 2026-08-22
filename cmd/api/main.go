package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_pgx_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/kupr666/Orange_Team/internal/core/transport/http/middleware"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

func main() {
	loggerConfig, err := core_logger.NewConfig()
	if err != nil {
		fmt.Println("failed to get logger config:", err)
		os.Exit(1)
	}

	log, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		log.Error("failed to init postgres connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	/*

		WORKOUT FEATURE

	*/

	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		log,
		core_http_middleware.RequestID(),
	)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", "error", err)
	}
}

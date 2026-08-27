package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_auth_jwt "github.com/kupr666/Orange_Team/internal/core/auth/jwt"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_pgx_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/kupr666/Orange_Team/internal/core/transport/http/middleware"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
	authentication_postgres_repository "github.com/kupr666/Orange_Team/internal/features/authentication/repository/postgres"
	authentication_service "github.com/kupr666/Orange_Team/internal/features/authentication/service"
	authentication_transport_http "github.com/kupr666/Orange_Team/internal/features/authentication/transport/http"
	exercises_postgres_repository "github.com/kupr666/Orange_Team/internal/features/exercises/repository/postgres"
	exercises_service "github.com/kupr666/Orange_Team/internal/features/exercises/service"
	exercises_transport_http "github.com/kupr666/Orange_Team/internal/features/exercises/transport/http"
	workouts_postgres_repository "github.com/kupr666/Orange_Team/internal/features/workouts/repository/postgres"
	workouts_service "github.com/kupr666/Orange_Team/internal/features/workouts/service"
	workouts_transport_http "github.com/kupr666/Orange_Team/internal/features/workouts/transport/http"
)

func main() {
	log, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer log.Close()

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

	log.Debug("initializing feature", "feature", "exercises")
	exercisesRepository := exercises_postgres_repository.NewExercisesRepository(pool)
	exercisesService := exercises_service.NewExercisesService(exercisesRepository)
	exercisesTransportHTTP := exercises_transport_http.NewExercisesHTTPHandler(exercisesService)

	log.Debug("initializing feature", "feature", "exercises")
	workoutsRepository := workouts_postgres_repository.NewWorkoutsRepository(pool)
	workoutsService := workouts_service.NewWorkoutsService(workoutsRepository)
	workoutsTransportHTTP := workouts_transport_http.NewWorkoutsHTTPHandler(workoutsService)

	log.Debug("initializing feature", "feature", "authentication")
	authenticationRepository := authentication_postgres_repository.NewAuthenticationRepository(pool)
	jwtConfig := core_auth_jwt.NewConfigMust()
	jwtManager, err := core_auth_jwt.NewManager(jwtConfig)
	if err != nil {
		log.Error(
			"failed to initialize jwt manager",
			"error",
			err,
		)
		os.Exit(1)
	}
	authenticationService, err := authentication_service.NewAuthenticationService(authenticationRepository, jwtManager)
	if err != nil {
		log.Error("failed to initialize authentication service", "error", err)
		os.Exit(1)
	}
	authenticationTransportHTTP := authentication_transport_http.NewAuthenticationHTTPHandler(authenticationService)

	/*

		SOME FEATURE

	*/

	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		log,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(log),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	authenticationMiddleware := core_http_middleware.Authentication(jwtManager)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(authenticationTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(exercisesTransportHTTP.Routes(authenticationMiddleware, )...)
	apiVersionRouterV1.RegisterRoutes(workoutsTransportHTTP.Routes(authenticationMiddleware)...)

	// apiVersionRouterV1.RegisterRoutes(authenticationTransportHTTP.Routes()...)

	httpServer.RegisterRouters(
		apiVersionRouterV1,
	)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", "error", err)
	}
}

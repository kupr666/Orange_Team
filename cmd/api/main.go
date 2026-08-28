package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	api_docs "github.com/kupr666/Orange_Team/docs"
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
	leaderboard_postgres_repository "github.com/kupr666/Orange_Team/internal/features/leaderboard/repository/postgres"
	leaderboard_service "github.com/kupr666/Orange_Team/internal/features/leaderboard/service"
	leaderboard_transport_http "github.com/kupr666/Orange_Team/internal/features/leaderboard/transport/http"
	leaderboard_worker "github.com/kupr666/Orange_Team/internal/features/leaderboard/worker"
	users_postgres_repository "github.com/kupr666/Orange_Team/internal/features/users/repository/postgres"
	users_service "github.com/kupr666/Orange_Team/internal/features/users/service"
	users_transport_http "github.com/kupr666/Orange_Team/internal/features/users/transport/http"
	workout_exercises_postgres_repository "github.com/kupr666/Orange_Team/internal/features/workout_exercises/repository/postgres"
	workout_exercises_service "github.com/kupr666/Orange_Team/internal/features/workout_exercises/service"
	workout_exercises_transport_http "github.com/kupr666/Orange_Team/internal/features/workout_exercises/transport/http"
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

	log.Debug("initializing repositories")
	exercisesRepository := exercises_postgres_repository.NewExercisesRepository(pool)
	workoutsRepository := workouts_postgres_repository.NewWorkoutsRepository(pool)
	workoutExercisesRepo := workout_exercises_postgres_repository.NewWorkoutExercisesRepository(pool)
	authenticationRepository := authentication_postgres_repository.NewAuthenticationRepository(pool)
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	leaderboardRepository := leaderboard_postgres_repository.NewLeaderboardRepository(pool)

	// jwt manager (verifier)
	jwtConfig := core_auth_jwt.NewConfigMust()
	jwtManager, err := core_auth_jwt.NewManager(jwtConfig)
	if err != nil {
		log.Error("failed to initialize jwt manager", "error", err)
		os.Exit(1)
	}

	// leaderboard
	leaderboardConfig := leaderboard_service.NewConfigMust()

	log.Debug("initializing services")
	exercisesService := exercises_service.NewExercisesService(exercisesRepository)
	workoutsService := workouts_service.NewWorkoutsService(
		workoutsRepository,
		workoutExercisesRepo,
	)
	workoutExercisesService := workout_exercises_service.NewWorkoutExercisesService(
		workoutExercisesRepo,
		exercisesRepository,
		workoutsRepository,
		workoutExercisesRepo,
	)
	authenticationService, err := authentication_service.NewAuthenticationService(authenticationRepository, jwtManager)
	if err != nil {
		log.Error("failed to initialize authentication service", "error", err)
		os.Exit(1)
	}
	usersService := users_service.NewUsersService(usersRepository)
	leaderboardService := leaderboard_service.NewLeaderboardService(
		leaderboardRepository,
		leaderboardConfig.LocationMust(),
	)

	log.Debug("initializing HTTP handlers")
	exercisesTransportHTTP := exercises_transport_http.NewExercisesHTTPHandler(exercisesService)
	workoutsTransportHTTP := workouts_transport_http.NewWorkoutsHTTPHandler(workoutsService)
	workoutExercisesTransportHTTP := workout_exercises_transport_http.NewWorkoutExercisesHandler(workoutExercisesService)
	authenticationTransportHTTP := authentication_transport_http.NewAuthenticationHTTPHandler(authenticationService)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)
	leaderboardTransportHTTP := leaderboard_transport_http.NewLeaderboardHTTPHandler(leaderboardService)

	// leaderboard
	leaderboardSnapshotWorker := leaderboard_worker.NewSnapshotWorker(
		leaderboardService,
		log,
		leaderboardConfig.SnapshotInterval,
	)
	go leaderboardSnapshotWorker.Run(ctx)

	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		log,
		core_http_middleware.RequestID(),
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.Logger(log),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	authenticationMiddleware := core_http_middleware.Authentication(jwtManager)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(authenticationTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(exercisesTransportHTTP.Routes(authenticationMiddleware)...)
	apiVersionRouterV1.RegisterRoutes(workoutsTransportHTTP.Routes(authenticationMiddleware)...)
	apiVersionRouterV1.RegisterRoutes(leaderboardTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes(authenticationMiddleware)...)
	apiVersionRouterV1.RegisterRoutes(workoutExercisesTransportHTTP.Routes(authenticationMiddleware)...)

	httpServer.RegisterRouters(
		apiVersionRouterV1,
	)
	httpServer.RegisterRoutes(api_docs.Routes()...)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", "error", err)
	}
}
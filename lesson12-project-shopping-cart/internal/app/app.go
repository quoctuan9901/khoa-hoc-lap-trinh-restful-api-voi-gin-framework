package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/db"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/routes"
	"user-management-api/internal/validation"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"
	"user-management-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB    sqlc.Querier
	Redis *redis.Client
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Validator init failed")
	}

	r := gin.Default()

	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database init failed")
	}

	redisClient := config.NewRedisClient()
	cacheRedisService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedisService)

	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
	}

	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, cacheRedisService),
	}

	routes.RegisterRoutes(r, tokenService, cacheRedisService, getModulRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		logger.Log.Info().Msgf("✅ Server is running at %s", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("⛔️ Failed to start server")
		}
	}()

	<-quit
	logger.Log.Warn().Msg("⚠️  Shutdown signal received ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("⛔️ Server forced to shutdown")
	}

	logger.Log.Info().Msg("🍺 Server exited gracefully")

	return nil
}

func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}

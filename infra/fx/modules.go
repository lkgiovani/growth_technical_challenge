package fx

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lkgiovani/growth_technical_challenge/infra/cache"
	"github.com/lkgiovani/growth_technical_challenge/infra/config"
	"github.com/lkgiovani/growth_technical_challenge/infra/database"
	httpDelivery "github.com/lkgiovani/growth_technical_challenge/internal/delivery/http"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/httpError"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/middleware"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/router"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/colaborador"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/departamento"

	"github.com/lkgiovani/growth_technical_challenge/internal/repository"
	"github.com/lkgiovani/growth_technical_challenge/internal/services/cache_service"
	pkgCache "github.com/lkgiovani/growth_technical_challenge/pkg/cache"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		NewLogger,
		NewDatabase,
		NewRedisClient,
		NewCacheService,
		NewRouter,
		NewLoggerMiddleware,
		NewRecoveryMiddleware,
		NewColaboradorRepository,
		NewDepartamentoRepository,
		NewColaboradorUseCase,
		NewDepartamentoUseCase,
		httpError.NewErrorHandler,
		httpDelivery.NewColaboradorHandler,
		httpDelivery.NewDepartamentoHandler,
		httpDelivery.NewGerenteHandler,
		httpDelivery.NewDocsHandler,
		middleware.NewPrometheusMetrics,
	),
	fx.Invoke(RegisterLifecycle),
)

func NewLogger(lc fx.Lifecycle, cfg *config.Config) (logger.Logger, error) {
	log, err := logger.NewLoggerWithConfig(cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Logger initialized successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down logger...")
			return log.Sync()
		},
	})

	return log, nil
}

func NewDatabase(lc fx.Lifecycle, cfg *config.Config, log logger.Logger) (*gorm.DB, error) {
	if err := database.Connect(); err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	db := database.GetDB()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Database connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				log.Error("Failed to get database instance", zap.Error(err))
				return err
			}
			log.Info("Closing database connection...")
			return sqlDB.Close()
		},
	})

	return db, nil
}

func NewRedisClient(lc fx.Lifecycle, cfg *config.Config, log logger.Logger) (pkgCache.Cache, error) {
	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Warn("Failed to connect to Redis. Continuing without cache", zap.Error(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Redis connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing Redis connection...")
			return redisClient.Close()
		},
	})

	return redisClient, nil
}

type LoggerMiddleware gin.HandlerFunc
type RecoveryMiddleware gin.HandlerFunc

func NewRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	return gin.New()
}

func NewLoggerMiddleware(log logger.Logger) LoggerMiddleware {
	return LoggerMiddleware(middleware.Logger(log))
}

func NewRecoveryMiddleware(log logger.Logger) RecoveryMiddleware {
	return RecoveryMiddleware(middleware.Recovery(log))
}

func NewCacheService(cache pkgCache.Cache, log logger.Logger) cache_service.CacheService {
	if cache == nil {
		log.Warn("Cache not available, cache service will operate in pass-through mode")
	}
	return cache_service.NewCacheService(cache, log)
}

func NewColaboradorRepository(db *gorm.DB, log logger.Logger) repository.ColaboradorRepository {
	return repository.NewColaboradorRepository(db, log)
}

func NewDepartamentoRepository(db *gorm.DB, log logger.Logger) repository.DepartamentoRepository {
	return repository.NewDepartamentoRepository(db, log)
}

func NewColaboradorUseCase(
	colaboradorRepo repository.ColaboradorRepository,
	departamentoRepo repository.DepartamentoRepository,
	cacheService cache_service.CacheService,
	db *gorm.DB,
	log logger.Logger,
) colaborador.UseCase {
	return colaborador.NewUseCase(colaboradorRepo, departamentoRepo, cacheService, db, log)
}

func NewDepartamentoUseCase(
	departamentoRepo repository.DepartamentoRepository,
	colaboradorRepo repository.ColaboradorRepository,
	cacheService cache_service.CacheService,
	db *gorm.DB,
	log logger.Logger,
) departamento.UseCase {
	return departamento.NewUseCase(departamentoRepo, colaboradorRepo, cacheService, db, log)
}

type RouteParams struct {
	fx.In

	Config              *config.Config
	Router              *gin.Engine
	ColaboradorHandler  *httpDelivery.ColaboradorHandler
	DepartamentoHandler *httpDelivery.DepartamentoHandler
	GerenteHandler      *httpDelivery.GerenteHandler
	DocsHandler         *httpDelivery.DocsHandler
	PrometheusMetrics   *middleware.PrometheusMetrics
	LoggerMiddleware    LoggerMiddleware
	RecoveryMiddleware  RecoveryMiddleware
	Logger              logger.Logger
}

func RegisterLifecycle(lc fx.Lifecycle, p RouteParams) {
	router.SetupRoutes(p.Router, p.ColaboradorHandler, p.DepartamentoHandler, p.GerenteHandler, p.DocsHandler, p.PrometheusMetrics, gin.HandlerFunc(p.LoggerMiddleware), gin.HandlerFunc(p.RecoveryMiddleware))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := fmt.Sprintf(":%s", p.Config.Server.Port)
				p.Logger.Info("Starting server",
					zap.String("port", p.Config.Server.Port),
					zap.String("mode", p.Config.Server.Mode),
				)
				p.Logger.Info("API documentation available at:",
					zap.String("ReDoc", fmt.Sprintf("http://localhost:%s/docs/redoc", p.Config.Server.Port)),
					zap.String("Swagger", fmt.Sprintf("http://localhost:%s/docs/swagger", p.Config.Server.Port)),
					zap.String("Scalar", fmt.Sprintf("http://localhost:%s/docs/scalar", p.Config.Server.Port)),
					zap.String("OpenAPI", fmt.Sprintf("http://localhost:%s/docs/openapi.yaml", p.Config.Server.Port)),
				)
				p.Logger.Info("Prometheus metrics available at:",
					zap.String("Metrics", fmt.Sprintf("http://localhost:%s/metrics", p.Config.Server.Port)),
				)
				if err := p.Router.Run(addr); err != nil {
					p.Logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Shutting down server...")
			return nil
		},
	})
}

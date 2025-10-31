package fx

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lkgiovani/growth_technical_challenge/infra/cache"
	"github.com/lkgiovani/growth_technical_challenge/infra/config"
	"github.com/lkgiovani/growth_technical_challenge/infra/database"
	httpDelivery "github.com/lkgiovani/growth_technical_challenge/internal/delivery/http"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/router"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/colaborador"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/departamento"

	"github.com/lkgiovani/growth_technical_challenge/internal/repository"
	pkgCache "github.com/lkgiovani/growth_technical_challenge/pkg/cache"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		NewDatabase,
		NewRedisClient,
		NewRouter,
		NewColaboradorRepository,
		NewDepartamentoRepository,
		NewColaboradorUseCase,
		NewDepartamentoUseCase,
		httpDelivery.NewColaboradorHandler,
		httpDelivery.NewDepartamentoHandler,
		httpDelivery.NewGerenteHandler,
		httpDelivery.NewDocsHandler,
	),
	fx.Invoke(RegisterLifecycle),
)

func NewDatabase(lc fx.Lifecycle, cfg *config.Config) (*gorm.DB, error) {
	if err := database.Connect(); err != nil {
		return nil, err
	}

	db := database.GetDB()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Banco de dados conectado com sucesso")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			log.Println("Fechando conexão com o banco de dados...")
			return sqlDB.Close()
		},
	})

	return db, nil
}

func NewRedisClient(lc fx.Lifecycle, cfg *config.Config) (pkgCache.Cache, error) {
	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Printf("Aviso: Falha ao conectar ao Redis: %v. Continuando sem cache.", err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Redis conectado com sucesso")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Fechando conexão com Redis...")
			return redisClient.Close()
		},
	})

	return redisClient, nil
}

func NewRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	return gin.Default()
}

func NewColaboradorRepository(db *gorm.DB, cache pkgCache.Cache) repository.ColaboradorRepository {
	if cache == nil {
		log.Println("Cache não disponível para Colaborador, usando repositório sem cache")
	}
	return repository.NewColaboradorRepository(db, cache)
}

func NewDepartamentoRepository(db *gorm.DB, cache pkgCache.Cache) repository.DepartamentoRepository {
	if cache == nil {
		log.Println("Cache não disponível para Departamento, usando repositório sem cache")
	}
	return repository.NewDepartamentoRepository(db, cache)
}

func NewColaboradorUseCase(
	colaboradorRepo repository.ColaboradorRepository,
	departamentoRepo repository.DepartamentoRepository,
) colaborador.UseCase {
	return colaborador.NewUseCase(colaboradorRepo, departamentoRepo)
}

func NewDepartamentoUseCase(
	departamentoRepo repository.DepartamentoRepository,
	colaboradorRepo repository.ColaboradorRepository,
	db *gorm.DB,
) departamento.UseCase {
	return departamento.NewUseCase(departamentoRepo, colaboradorRepo, db)
}

type RouteParams struct {
	fx.In

	Config              *config.Config
	Router              *gin.Engine
	ColaboradorHandler  *httpDelivery.ColaboradorHandler
	DepartamentoHandler *httpDelivery.DepartamentoHandler
	GerenteHandler      *httpDelivery.GerenteHandler
	DocsHandler         *httpDelivery.DocsHandler
}

func RegisterLifecycle(lc fx.Lifecycle, p RouteParams) {
	router.SetupRoutes(p.Router, p.ColaboradorHandler, p.DepartamentoHandler, p.GerenteHandler, p.DocsHandler)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := fmt.Sprintf(":%s", p.Config.Server.Port)
				log.Printf("Iniciando servidor na porta %s", p.Config.Server.Port)
				log.Printf("Documentação da API disponível em:")
				log.Printf("  - ReDoc:        http://localhost:%s/docs/redoc", p.Config.Server.Port)
				log.Printf("  - Swagger UI:   http://localhost:%s/docs/swagger", p.Config.Server.Port)
				log.Printf("  - Scalar:       http://localhost:%s/docs/scalar", p.Config.Server.Port)
				log.Printf("  - OpenAPI Spec: http://localhost:%s/docs/openapi.yaml", p.Config.Server.Port)
				if err := p.Router.Run(addr); err != nil {
					log.Fatalf("Falha ao iniciar servidor: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Desligando servidor...")
			return nil
		},
	})
}

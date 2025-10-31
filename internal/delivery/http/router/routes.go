package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpDelivery "github.com/lkgiovani/growth_technical_challenge/internal/delivery/http"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(
	r *gin.Engine,
	colaboradorHandler *httpDelivery.ColaboradorHandler,
	departamentoHandler *httpDelivery.DepartamentoHandler,
	gerenteHandler *httpDelivery.GerenteHandler,
	docsHandler *httpDelivery.DocsHandler,
	prometheusMetrics *middleware.PrometheusMetrics,
	loggerMiddleware gin.HandlerFunc,
	recoveryMiddleware gin.HandlerFunc,
) {
	r.Use(recoveryMiddleware)
	r.Use(loggerMiddleware)
	r.Use(prometheusMetrics.Middleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.LoadHTMLGlob("docs/schema/*.html")

	r.GET("/favicon.ico", httpDelivery.ServeFavicon)
	r.GET("/", httpDelivery.ServeIndex)

	docs := r.Group("/docs")
	{
		docs.GET("", docsHandler.DocsIndex)
		docs.GET("/openapi.yaml", docsHandler.ServeOpenAPI)
		docs.GET("/redoc", docsHandler.ServeRedoc)
		docs.GET("/swagger", docsHandler.ServeSwagger)
		docs.GET("/scalar", docsHandler.ServeScalar)
	}

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":   "ok",
				"database": "connected",
			})
		})

		colaboradores := api.Group("/colaboradores")
		{
			colaboradores.POST("", colaboradorHandler.CreateColaborador)
			colaboradores.GET("/:id", colaboradorHandler.GetColaborador)
			colaboradores.PUT("/:id", colaboradorHandler.UpdateColaborador)
			colaboradores.DELETE("/:id", colaboradorHandler.DeleteColaborador)
			colaboradores.POST("/listar", colaboradorHandler.ListColaboradores)
		}

		departamentos := api.Group("/departamentos")
		{
			departamentos.POST("", departamentoHandler.CreateDepartamento)
			departamentos.GET("/:id", departamentoHandler.GetDepartamento)
			departamentos.PUT("/:id", departamentoHandler.UpdateDepartamento)
			departamentos.DELETE("/:id", departamentoHandler.DeleteDepartamento)
			departamentos.POST("/listar", departamentoHandler.ListDepartamentos)
		}

		gerentes := api.Group("/gerentes")
		{
			gerentes.GET("/:id/colaboradores", gerenteHandler.GetGerenteColaboradores)
		}
	}
}

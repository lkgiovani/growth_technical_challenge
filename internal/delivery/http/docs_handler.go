package http

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) ServeOpenAPI(c *gin.Context) {
	schemaPath := filepath.Join("docs", "schema", "openapi.yaml")
	c.File(schemaPath)
}

func (h *DocsHandler) ServeRedoc(c *gin.Context) {
	c.HTML(http.StatusOK, "redoc.html", nil)
}

func (h *DocsHandler) ServeSwagger(c *gin.Context) {
	c.HTML(http.StatusOK, "swagger.html", nil)
}

func (h *DocsHandler) ServeScalar(c *gin.Context) {
	c.HTML(http.StatusOK, "scalar.html", nil)
}

func (h *DocsHandler) DocsIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "API Documentation",
		"viewers": gin.H{
			"openapi": "/docs/openapi.yaml",
			"redoc":   "/docs/redoc",
			"swagger": "/docs/swagger",
			"scalar":  "/docs/scalar",
		},
	})
}

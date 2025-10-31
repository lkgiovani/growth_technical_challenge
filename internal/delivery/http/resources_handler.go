package http

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed resources/index.html
var indexHTML string

//go:embed resources/grownt.ico
var faviconICO []byte

func ServeIndex(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, indexHTML)
}

func ServeFavicon(c *gin.Context) {
	c.Data(http.StatusOK, "image/x-icon", faviconICO)
}

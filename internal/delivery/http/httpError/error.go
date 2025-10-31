package httpError

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	app "github.com/lkgiovani/growth_technical_challenge"
)

func NewErrorHandler() func(c *gin.Context, err error) {
	return Error
}

var codes = map[string]int{
	app.ECONFLICT:       http.StatusConflict,
	app.EINVALID:        http.StatusBadRequest,
	app.ENOTFOUND:       http.StatusNotFound,
	app.ENOTIMPLEMENTED: http.StatusNotImplemented,
	app.EUNAUTHORIZED:   http.StatusUnauthorized,
	app.EINTERNAL:       http.StatusInternalServerError,
	app.EDUPLICATION:    http.StatusConflict,
	app.EBADREQUEST:     http.StatusBadRequest,
	app.EFORBIDDEN:      http.StatusForbidden,
	app.ETIMEOUT:        http.StatusRequestTimeout,
	app.EUNAVAILABLE:    http.StatusServiceUnavailable,
}

func Error(c *gin.Context, err error) {
	code, message := app.ErrorCode(err), app.ErrorMessage(err)
	if code == app.EINTERNAL {
		LogError(c, err)
	}

	c.JSON(ErrorStatusCode(code), gin.H{"error": message})
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

func LogError(c *gin.Context, err error) {
	log.Printf("[http] error: %s %s: %s ", c.Request.Method, c.Request.URL.Path, err)
}

package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/departamento"
)

type GerenteHandler struct {
	departamentoUseCase departamento.UseCase
}

func NewGerenteHandler(departamentoUseCase departamento.UseCase) *GerenteHandler {
	return &GerenteHandler{
		departamentoUseCase: departamentoUseCase,
	}
}

func (h *GerenteHandler) GetGerenteColaboradores(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid manager ID",
			Message: "ID must be a valid UUID",
		})
		return
	}

	colaboradores, err := h.departamentoUseCase.GetGerenteColaboradores(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch app.ErrorCode(err) {
		case app.ENOTFOUND:
			statusCode = http.StatusNotFound
		case app.EINVALID:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to retrieve manager employees",
			Message: app.ErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data:  colaboradores,
		Count: len(colaboradores),
	})
}

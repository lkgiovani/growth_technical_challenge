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
			Error:   "ID de gerente inválido",
			Message: "ID deve ser um UUID válido",
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
			Error:   "Falha ao recuperar colaboradores do gerente",
			Message: app.ErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data:  colaboradores,
		Count: len(colaboradores),
	})
}

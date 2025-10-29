package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/usecases/departamento"
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
		if err.Error() == "gerente não encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Falha ao recuperar colaboradores do gerente",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data:  colaboradores,
		Count: len(colaboradores),
	})
}

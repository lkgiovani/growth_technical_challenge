package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/departamento"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"go.uber.org/zap"
)

type GerenteHandler struct {
	departamentoUseCase departamento.UseCase
	errorHandler        func(c *gin.Context, err error)
	logger              logger.Logger
}

func NewGerenteHandler(departamentoUseCase departamento.UseCase, errorHandler func(c *gin.Context, err error), log logger.Logger) *GerenteHandler {
	return &GerenteHandler{
		departamentoUseCase: departamentoUseCase,
		errorHandler:        errorHandler,
		logger:              log,
	}
}

func (h *GerenteHandler) GetGerenteColaboradores(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "ID must be a valid UUID"))
		return
	}

	h.logger.Info("Fetching colaboradores for manager", zap.String("manager_id", id.String()))

	colaboradores, err := h.departamentoUseCase.GetGerenteColaboradores(id)
	if err != nil {
		h.logger.Error("Failed to fetch colaboradores for manager", zap.Error(err), zap.String("manager_id", id.String()))
		h.errorHandler(c, err)
		return
	}

	h.logger.Info("Fetched colaboradores for manager successfully", zap.String("manager_id", id.String()), zap.Int("count", len(colaboradores)))

	c.JSON(http.StatusOK, SuccessResponse{
		Data:  colaboradores,
		Count: len(colaboradores),
	})
}

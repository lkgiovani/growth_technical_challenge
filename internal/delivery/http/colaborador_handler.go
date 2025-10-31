package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/dto"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/colaborador"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"go.uber.org/zap"
)

type ColaboradorHandler struct {
	colaboradorUseCase colaborador.UseCase
	errorHandler       func(c *gin.Context, err error)
	logger             logger.Logger
}

func NewColaboradorHandler(colaboradorUseCase colaborador.UseCase, errorHandler func(c *gin.Context, err error), log logger.Logger) *ColaboradorHandler {
	return &ColaboradorHandler{
		colaboradorUseCase: colaboradorUseCase,
		errorHandler:       errorHandler,
		logger:             log,
	}
}

func (h *ColaboradorHandler) GetColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "ID must be a valid UUID"))
		return
	}

	colaborador, err := h.colaboradorUseCase.GetColaboradorByID(id)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: colaborador,
	})
}

func (h *ColaboradorHandler) CreateColaborador(c *gin.Context) {
	var req dto.CreateColaboradorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "Invalid request body: %v", err))
		return
	}

	h.logger.Info("Creating colaborador", zap.String("name", req.Nome), zap.String("cpf", req.CPF))

	colaborador := &entities.Colaborador{
		Nome:           req.Nome,
		CPF:            req.CPF,
		RG:             req.RG,
		DepartamentoID: req.DepartamentoID,
	}

	if err := h.colaboradorUseCase.CreateColaborador(colaborador); err != nil {
		h.logger.Error("Failed to create colaborador", zap.Error(err), zap.String("cpf", req.CPF))
		h.errorHandler(c, err)
		return
	}

	h.logger.Info("Colaborador created successfully", zap.String("id", colaborador.ID.String()), zap.String("name", colaborador.Nome))

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Employee created successfully",
		Data:    colaborador,
	})
}

func (h *ColaboradorHandler) UpdateColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "ID must be a valid UUID"))
		return
	}

	h.logger.Info("Updating colaborador", zap.String("id", id.String()))

	var req dto.UpdateColaboradorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "Invalid request body: %v", err))
		return
	}

	colaborador := &entities.Colaborador{
		ID: id,
	}

	if req.Nome != nil {
		colaborador.Nome = *req.Nome
	}

	if req.CPF != nil {
		colaborador.CPF = *req.CPF
	}

	if req.RG != nil {
		colaborador.RG = req.RG
	}

	if req.DepartamentoID != nil {
		colaborador.DepartamentoID = *req.DepartamentoID
	}

	if err := h.colaboradorUseCase.UpdateColaborador(id, colaborador); err != nil {
		h.logger.Error("Failed to update colaborador", zap.Error(err), zap.String("id", id.String()))
		h.errorHandler(c, err)
		return
	}

	h.logger.Info("Colaborador updated successfully", zap.String("id", id.String()))

	updatedColaborador, _ := h.colaboradorUseCase.GetColaboradorByID(id)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Employee updated successfully",
		Data:    updatedColaborador,
	})
}

func (h *ColaboradorHandler) DeleteColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "ID must be a valid UUID"))
		return
	}

	h.logger.Info("Deleting colaborador", zap.String("id", id.String()))

	if err := h.colaboradorUseCase.DeleteColaborador(id); err != nil {
		h.logger.Error("Failed to delete colaborador", zap.Error(err), zap.String("id", id.String()))
		h.errorHandler(c, err)
		return
	}

	h.logger.Info("Colaborador deleted successfully", zap.String("id", id.String()))

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Employee deleted successfully",
	})
}

func (h *ColaboradorHandler) ListColaboradores(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "Invalid request body: %v", err))
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	colaboradores, total, err := h.colaboradorUseCase.ListColaboradores(req.Filters, req.Limit, offset)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       colaboradores,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	})
}

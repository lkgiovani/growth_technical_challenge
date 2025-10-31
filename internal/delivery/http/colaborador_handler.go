package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/dto"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/colaborador"
)

type ColaboradorHandler struct {
	colaboradorUseCase colaborador.UseCase
}

func NewColaboradorHandler(colaboradorUseCase colaborador.UseCase) *ColaboradorHandler {
	return &ColaboradorHandler{
		colaboradorUseCase: colaboradorUseCase,
	}
}

func (h *ColaboradorHandler) GetColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid employee ID",
			Message: "ID must be a valid UUID",
		})
		return
	}

	colaborador, err := h.colaboradorUseCase.GetColaboradorByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Employee not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: colaborador,
	})
}

func (h *ColaboradorHandler) CreateColaborador(c *gin.Context) {
	var req dto.CreateColaboradorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	colaborador := &entities.Colaborador{
		Nome:           req.Nome,
		CPF:            req.CPF,
		RG:             req.RG,
		DepartamentoID: req.DepartamentoID,
	}

	if err := h.colaboradorUseCase.CreateColaborador(colaborador); err != nil {
		statusCode := http.StatusUnprocessableEntity
		switch app.ErrorCode(err) {
		case app.EDUPLICATION:
			statusCode = http.StatusConflict
		case app.ENOTFOUND:
			statusCode = http.StatusNotFound
		case app.EINVALID:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to create employee",
			Message: app.ErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Employee created successfully",
		Data:    colaborador,
	})
}

func (h *ColaboradorHandler) UpdateColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid employee ID",
			Message: "ID must be a valid UUID",
		})
		return
	}

	var req dto.UpdateColaboradorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
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
		statusCode := http.StatusUnprocessableEntity
		switch app.ErrorCode(err) {
		case app.ENOTFOUND:
			statusCode = http.StatusNotFound
		case app.EDUPLICATION:
			statusCode = http.StatusConflict
		case app.EINVALID:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to update employee",
			Message: app.ErrorMessage(err),
		})
		return
	}

	updatedColaborador, _ := h.colaboradorUseCase.GetColaboradorByID(id)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Employee updated successfully",
		Data:    updatedColaborador,
	})
}

func (h *ColaboradorHandler) DeleteColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid employee ID",
			Message: "ID must be a valid UUID",
		})
		return
	}

	if err := h.colaboradorUseCase.DeleteColaborador(id); err != nil {
		statusCode := http.StatusInternalServerError
		switch app.ErrorCode(err) {
		case app.ENOTFOUND:
			statusCode = http.StatusNotFound
		case app.EINVALID:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to delete employee",
			Message: app.ErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Employee deleted successfully",
	})
}

func (h *ColaboradorHandler) ListColaboradores(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list employees",
			Message: err.Error(),
		})
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

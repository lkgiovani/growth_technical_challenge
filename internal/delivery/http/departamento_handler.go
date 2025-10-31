package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/dto"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/departamento"
)

type DepartamentoHandler struct {
	departamentoUseCase departamento.UseCase
}

func NewDepartamentoHandler(departamentoUseCase departamento.UseCase) *DepartamentoHandler {
	return &DepartamentoHandler{
		departamentoUseCase: departamentoUseCase,
	}
}

func (h *DepartamentoHandler) GetDepartamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid department ID",
			Message: "ID must be a valid UUID",
		})
		return
	}

	departamento, err := h.departamentoUseCase.GetDepartamentoWithHierarchy(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Department not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: departamento,
	})
}

func (h *DepartamentoHandler) CreateDepartamento(c *gin.Context) {
	var req dto.CreateDepartamentoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	if req.GerenteID != nil && req.Gerente != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: "Provide gerente_id OR gerente, not both",
		})
		return
	}

	if req.GerenteID == nil && req.Gerente == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: "Provide gerente_id or gerente",
		})
		return
	}

	if req.Gerente != nil {
		departamento := &entities.Departamento{
			Nome:                   req.Nome,
			DepartamentoSuperiorID: req.DepartamentoSuperiorID,
		}

		colaborador := &entities.Colaborador{
			Nome: req.Gerente.Nome,
			CPF:  req.Gerente.CPF,
			RG:   req.Gerente.RG,
		}

		if err := h.departamentoUseCase.CreateDepartamentoWithGerente(departamento, colaborador); err != nil {
			statusCode := http.StatusUnprocessableEntity
			switch app.ErrorCode(err) {
			case app.ENOTFOUND:
				statusCode = http.StatusNotFound
			case app.ECONFLICT:
				statusCode = http.StatusConflict
			case app.EINVALID:
				statusCode = http.StatusBadRequest
			case app.EDUPLICATION:
				statusCode = http.StatusConflict
			}
			c.JSON(statusCode, ErrorResponse{
				Error:   "Failed to create department with manager",
				Message: app.ErrorMessage(err),
			})
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse{
			Message: "Department and manager created successfully",
			Data:    departamento,
		})
		return
	}

	departamento := &entities.Departamento{
		Nome:                   req.Nome,
		GerenteID:              *req.GerenteID,
		DepartamentoSuperiorID: req.DepartamentoSuperiorID,
	}

	if err := h.departamentoUseCase.CreateDepartamento(departamento); err != nil {
		statusCode := http.StatusUnprocessableEntity
		switch app.ErrorCode(err) {
		case app.ENOTFOUND:
			statusCode = http.StatusNotFound
		case app.ECONFLICT:
			statusCode = http.StatusConflict
		case app.EINVALID:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to create department",
			Message: app.ErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Department created successfully",
		Data:    departamento,
	})
}

func (h *DepartamentoHandler) UpdateDepartamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid department ID",
			Message: "ID must be a valid UUID",
		})
		return
	}

	var req dto.UpdateDepartamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	departamento := &entities.Departamento{
		ID: id,
	}

	if req.Nome != nil {
		departamento.Nome = *req.Nome
	}

	if req.GerenteID != nil {
		departamento.GerenteID = *req.GerenteID
	}

	if req.DepartamentoSuperiorID != nil {
		departamento.DepartamentoSuperiorID = req.DepartamentoSuperiorID
	}

	if err := h.departamentoUseCase.UpdateDepartamento(id, departamento); err != nil {
		statusCode := http.StatusUnprocessableEntity
		switch app.ErrorCode(err) {
		case app.ENOTFOUND:
			statusCode = http.StatusNotFound
		case app.ECONFLICT:
			statusCode = http.StatusConflict
		case app.EINVALID:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to update department",
			Message: app.ErrorMessage(err),
		})
		return
	}

	updatedDepartamento, _ := h.departamentoUseCase.GetDepartamentoByID(id)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Department updated successfully",
		Data:    updatedDepartamento,
	})
}

func (h *DepartamentoHandler) DeleteDepartamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid department ID",
			Message: "ID must be a valid UUID",
		})
		return
	}

	if err := h.departamentoUseCase.DeleteDepartamento(id); err != nil {
		statusCode := http.StatusInternalServerError
		switch app.ErrorCode(err) {
		case app.ENOTFOUND:
			statusCode = http.StatusNotFound
		case app.EINVALID:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to delete department",
			Message: app.ErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Department deleted successfully",
	})
}

func (h *DepartamentoHandler) ListDepartamentos(c *gin.Context) {
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

	departamentos, total, err := h.departamentoUseCase.ListDepartamentos(req.Filters, req.Limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list departments",
			Message: err.Error(),
		})
		return
	}

	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       departamentos,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	})
}

package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/delivery/http/dto"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/usecases/departamento"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"go.uber.org/zap"
)

type DepartamentoHandler struct {
	departamentoUseCase departamento.UseCase
	errorHandler        func(c *gin.Context, err error)
	logger              logger.Logger
}

func NewDepartamentoHandler(departamentoUseCase departamento.UseCase, errorHandler func(c *gin.Context, err error), log logger.Logger) *DepartamentoHandler {
	return &DepartamentoHandler{
		departamentoUseCase: departamentoUseCase,
		errorHandler:        errorHandler,
		logger:              log,
	}
}

func (h *DepartamentoHandler) GetDepartamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "ID must be a valid UUID"))
		return
	}

	departamento, err := h.departamentoUseCase.GetDepartamentoWithHierarchy(id)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: departamento,
	})
}

func (h *DepartamentoHandler) CreateDepartamento(c *gin.Context) {
	var req dto.CreateDepartamentoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "Invalid request body: %v", err))
		return
	}

	if req.GerenteID != nil && req.Gerente != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "Provide gerente_id OR gerente, not both"))
		return
	}

	if req.GerenteID == nil && req.Gerente == nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "Provide gerente_id or gerente"))
		return
	}

	if req.Gerente != nil {
		h.logger.Info("Creating department with new manager", zap.String("department_name", req.Nome), zap.String("manager_name", req.Gerente.Nome))

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
			h.logger.Error("Failed to create department with manager", zap.Error(err), zap.String("department_name", req.Nome))
			h.errorHandler(c, err)
			return
		}

		h.logger.Info("Department and manager created successfully", zap.String("department_id", departamento.ID.String()), zap.String("manager_id", colaborador.ID.String()))

		c.JSON(http.StatusCreated, SuccessResponse{
			Message: "Department and manager created successfully",
			Data:    departamento,
		})
		return
	}

	h.logger.Info("Creating department with existing manager", zap.String("department_name", req.Nome), zap.String("manager_id", req.GerenteID.String()))

	departamento := &entities.Departamento{
		Nome:                   req.Nome,
		GerenteID:              *req.GerenteID,
		DepartamentoSuperiorID: req.DepartamentoSuperiorID,
	}

	if err := h.departamentoUseCase.CreateDepartamento(departamento); err != nil {
		h.logger.Error("Failed to create department", zap.Error(err), zap.String("department_name", req.Nome))
		h.errorHandler(c, err)
		return
	}

	h.logger.Info("Department created successfully", zap.String("department_id", departamento.ID.String()))

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Department created successfully",
		Data:    departamento,
	})
}

func (h *DepartamentoHandler) UpdateDepartamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "ID must be a valid UUID"))
		return
	}

	h.logger.Info("Updating department", zap.String("id", id.String()))

	var req dto.UpdateDepartamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "Invalid request body: %v", err))
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
		h.logger.Error("Failed to update department", zap.Error(err), zap.String("id", id.String()))
		h.errorHandler(c, err)
		return
	}

	h.logger.Info("Department updated successfully", zap.String("id", id.String()))

	updatedDepartamento, _ := h.departamentoUseCase.GetDepartamentoByID(id)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Department updated successfully",
		Data:    updatedDepartamento,
	})
}

func (h *DepartamentoHandler) DeleteDepartamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.errorHandler(c, app.Errorf(app.EINVALID, "ID must be a valid UUID"))
		return
	}

	h.logger.Info("Deleting department", zap.String("id", id.String()))

	if err := h.departamentoUseCase.DeleteDepartamento(id); err != nil {
		h.logger.Error("Failed to delete department", zap.Error(err), zap.String("id", id.String()))
		h.errorHandler(c, err)
		return
	}

	h.logger.Info("Department deleted successfully", zap.String("id", id.String()))

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Department deleted successfully",
	})
}

func (h *DepartamentoHandler) ListDepartamentos(c *gin.Context) {
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

	departamentos, total, err := h.departamentoUseCase.ListDepartamentos(req.Filters, req.Limit, offset)
	if err != nil {
		h.errorHandler(c, err)
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

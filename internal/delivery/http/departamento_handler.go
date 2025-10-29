package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/usecases/departamento"
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
			Error:   "ID de departamento inválido",
			Message: "ID deve ser um UUID válido",
		})
		return
	}

	departamento, err := h.departamentoUseCase.GetDepartamentoWithHierarchy(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Departamento não encontrado",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: departamento,
	})
}

func (h *DepartamentoHandler) CreateDepartamento(c *gin.Context) {
	var departamento entities.Departamento

	if err := c.ShouldBindJSON(&departamento); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Corpo da requisição inválido",
			Message: err.Error(),
		})
		return
	}

	if err := h.departamentoUseCase.CreateDepartamento(&departamento); err != nil {
		statusCode := http.StatusUnprocessableEntity
		if err.Error() == "gerente não encontrado" || err.Error() == "departamento superior não encontrado" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "não é possível criar ciclo na hierarquia de departamentos" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Falha ao criar departamento",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Departamento criado com sucesso",
		Data:    departamento,
	})
}

func (h *DepartamentoHandler) UpdateDepartamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID de departamento inválido",
			Message: "ID deve ser um UUID válido",
		})
		return
	}

	var departamento entities.Departamento
	if err := c.ShouldBindJSON(&departamento); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Corpo da requisição inválido",
			Message: err.Error(),
		})
		return
	}

	if err := h.departamentoUseCase.UpdateDepartamento(id, &departamento); err != nil {
		statusCode := http.StatusUnprocessableEntity
		if err.Error() == "departamento não encontrado" || err.Error() == "gerente não encontrado" || err.Error() == "departamento superior não encontrado" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "não é possível criar ciclo na hierarquia de departamentos" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Falha ao atualizar departamento",
			Message: err.Error(),
		})
		return
	}

	updatedDepartamento, _ := h.departamentoUseCase.GetDepartamentoByID(id)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Departamento atualizado com sucesso",
		Data:    updatedDepartamento,
	})
}

func (h *DepartamentoHandler) DeleteDepartamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID de departamento inválido",
			Message: "ID deve ser um UUID válido",
		})
		return
	}

	if err := h.departamentoUseCase.DeleteDepartamento(id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "departamento não encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Falha ao remover departamento",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Departamento removido com sucesso",
	})
}

func (h *DepartamentoHandler) ListDepartamentos(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Corpo da requisição inválido",
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
			Error:   "Falha ao listar departamentos",
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


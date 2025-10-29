package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/usecases/colaborador"
)

type ColaboradorHandler struct {
	colaboradorUseCase colaborador.UseCase
}

func NewColaboradorHandler(colaboradorUseCase colaborador.UseCase) *ColaboradorHandler {
	return &ColaboradorHandler{
		colaboradorUseCase: colaboradorUseCase,
	}
}

type ListRequest struct {
	Filters map[string]interface{} `json:"filters"`
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func (h *ColaboradorHandler) GetColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID de colaborador inválido",
			Message: "ID deve ser um UUID válido",
		})
		return
	}

	colaborador, err := h.colaboradorUseCase.GetColaboradorByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Colaborador não encontrado",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: colaborador,
	})
}

func (h *ColaboradorHandler) CreateColaborador(c *gin.Context) {
	var colaborador entities.Colaborador

	if err := c.ShouldBindJSON(&colaborador); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Corpo da requisição inválido",
			Message: err.Error(),
		})
		return
	}

	if err := h.colaboradorUseCase.CreateColaborador(&colaborador); err != nil {
		statusCode := http.StatusUnprocessableEntity
		if err.Error() == "CPF já existe" || err.Error() == "RG já existe" {
			statusCode = http.StatusConflict
		} else if err.Error() == "departamento não encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Falha ao criar colaborador",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Colaborador criado com sucesso",
		Data:    colaborador,
	})
}

func (h *ColaboradorHandler) UpdateColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID de colaborador inválido",
			Message: "ID deve ser um UUID válido",
		})
		return
	}

	var colaborador entities.Colaborador
	if err := c.ShouldBindJSON(&colaborador); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Corpo da requisição inválido",
			Message: err.Error(),
		})
		return
	}

	if err := h.colaboradorUseCase.UpdateColaborador(id, &colaborador); err != nil {
		statusCode := http.StatusUnprocessableEntity
		if err.Error() == "colaborador não encontrado" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "CPF já existe" || err.Error() == "RG já existe" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Falha ao atualizar colaborador",
			Message: err.Error(),
		})
		return
	}

	updatedColaborador, _ := h.colaboradorUseCase.GetColaboradorByID(id)

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Colaborador atualizado com sucesso",
		Data:    updatedColaborador,
	})
}

func (h *ColaboradorHandler) DeleteColaborador(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID de colaborador inválido",
			Message: "ID deve ser um UUID válido",
		})
		return
	}

	if err := h.colaboradorUseCase.DeleteColaborador(id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "colaborador não encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{
			Error:   "Falha ao remover colaborador",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Colaborador removido com sucesso",
	})
}

func (h *ColaboradorHandler) ListColaboradores(c *gin.Context) {
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

	colaboradores, total, err := h.colaboradorUseCase.ListColaboradores(req.Filters, req.Limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Falha ao listar colaboradores",
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


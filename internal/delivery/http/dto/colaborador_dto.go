package dto

import (
	"github.com/google/uuid"
)

type CreateColaboradorRequest struct {
	Nome           string    `json:"nome" binding:"required,min=3,max=255"`
	CPF            string    `json:"cpf" binding:"required"`
	RG             *string   `json:"rg"`
	DepartamentoID uuid.UUID `json:"departamento_id" binding:"required"`
}

type UpdateColaboradorRequest struct {
	Nome           *string    `json:"nome" binding:"omitempty,min=3,max=255"`
	CPF            *string    `json:"cpf"`
	RG             *string    `json:"rg"`
	DepartamentoID *uuid.UUID `json:"departamento_id"`
}

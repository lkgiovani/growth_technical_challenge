package dto

import (
	"github.com/google/uuid"
)

type CreateGerenteRequest struct {
	Nome string  `json:"nome" binding:"required,min=3,max=255"`
	CPF  string  `json:"cpf" binding:"required"`
	RG   *string `json:"rg"`
}

type CreateDepartamentoRequest struct {
	Nome                   string                `json:"nome" binding:"required,min=3,max=255"`
	GerenteID              *uuid.UUID            `json:"gerente_id"`
	DepartamentoSuperiorID *uuid.UUID            `json:"departamento_superior_id"`
	Gerente                *CreateGerenteRequest `json:"gerente"`
}

type UpdateDepartamentoRequest struct {
	Nome                   *string    `json:"nome" binding:"omitempty,min=3,max=255"`
	GerenteID              *uuid.UUID `json:"gerente_id"`
	DepartamentoSuperiorID *uuid.UUID `json:"departamento_superior_id"`
}

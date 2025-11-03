package dto

import (
	"github.com/google/uuid"
)

type DepartamentoForColaborador struct {
	Nome                   string     `json:"nome" binding:"required,min=3,max=255"`
	GerenteID              uuid.UUID  `json:"gerente_id" binding:"required"`
	DepartamentoSuperiorID *uuid.UUID `json:"departamento_superior_id"`
}

type CreateColaboradorRequest struct {
	Nome           string                      `json:"nome" binding:"required,min=3,max=255"`
	CPF            string                      `json:"cpf" binding:"required"`
	RG             *string                     `json:"rg"`
	DepartamentoID *uuid.UUID                  `json:"departamento_id"`
	Departamento   *DepartamentoForColaborador `json:"departamento"`
}

type UpdateColaboradorRequest struct {
	Nome           *string    `json:"nome" binding:"omitempty,min=3,max=255"`
	CPF            *string    `json:"cpf"`
	RG             *string    `json:"rg"`
	DepartamentoID *uuid.UUID `json:"departamento_id"`
}

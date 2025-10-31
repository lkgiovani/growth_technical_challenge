package departamento

import (
	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
)

type UseCase interface {
	GetAllDepartamentos(limit, offset int) ([]entities.Departamento, int64, error)
	GetDepartamentoByID(id uuid.UUID) (*entities.Departamento, error)
	GetDepartamentoWithHierarchy(id uuid.UUID) (*entities.Departamento, error)
	CreateDepartamento(departamento *entities.Departamento) error
	CreateDepartamentoWithGerente(departamento *entities.Departamento, gerente *entities.Colaborador) error
	UpdateDepartamento(id uuid.UUID, departamento *entities.Departamento) error
	DeleteDepartamento(id uuid.UUID) error
	ListDepartamentos(filters map[string]interface{}, limit, offset int) ([]entities.Departamento, int64, error)
	GetGerenteColaboradores(gerenteID uuid.UUID) ([]entities.Colaborador, error)
}

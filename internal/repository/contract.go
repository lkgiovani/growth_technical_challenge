package repository

import (
	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
)

type ColaboradorRepository interface {
	FindAll(limit, offset int) ([]entities.Colaborador, int64, error)
	FindByID(id uuid.UUID) (*entities.Colaborador, error)
	FindByCPF(cpf string) (*entities.Colaborador, error)
	FindByRG(rg string) (*entities.Colaborador, error)
	FindByFilters(filters map[string]interface{}, limit, offset int) ([]entities.Colaborador, int64, error)
	FindByDepartmentIDs(departmentIDs []uuid.UUID) ([]entities.Colaborador, error)
	Create(colaborador *entities.Colaborador) error
	Update(colaborador *entities.Colaborador) error
	Delete(id uuid.UUID) error
	Count() (int64, error)
}

type DepartamentoRepository interface {
	FindAll(limit, offset int) ([]entities.Departamento, int64, error)
	FindByID(id uuid.UUID) (*entities.Departamento, error)
	FindByIDWithHierarchy(id uuid.UUID) (*entities.Departamento, error)
	FindByFilters(filters map[string]interface{}, limit, offset int) ([]entities.Departamento, int64, error)
	FindSubDepartments(parentID uuid.UUID) ([]entities.Departamento, error)
	FindAllSubDepartmentIDs(parentID uuid.UUID) ([]uuid.UUID, error)
	HasCycle(departmentID, parentID uuid.UUID) (bool, error)
	Create(departamento *entities.Departamento) error
	Update(departamento *entities.Departamento) error
	Delete(id uuid.UUID) error
	Count() (int64, error)
}

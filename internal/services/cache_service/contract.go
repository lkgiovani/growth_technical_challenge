package cache_service

import (
	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
)

type CacheService interface {
	GetColaboradoresByDepartment(departmentID uuid.UUID) ([]entities.Colaborador, error)
	SetColaboradoresByDepartment(departmentID uuid.UUID, colaboradores []entities.Colaborador) error
	InvalidateColaboradoresByDepartment(departmentID uuid.UUID) error

	GetDepartmentHierarchy(departmentID uuid.UUID) (*entities.Departamento, error)
	SetDepartmentHierarchy(departmentID uuid.UUID, department *entities.Departamento) error
	InvalidateDepartmentHierarchy(departmentID uuid.UUID) error

	GetSubDepartmentIDs(parentID uuid.UUID) ([]uuid.UUID, error)
	SetSubDepartmentIDs(parentID uuid.UUID, ids []uuid.UUID) error
	InvalidateSubDepartmentIDs(parentID uuid.UUID) error

	InvalidateAllDepartmentCache(departmentID uuid.UUID, parentID *uuid.UUID) error
}

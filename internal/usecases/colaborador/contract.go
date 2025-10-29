package colaborador

import (
	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
)

type UseCase interface {
	GetAllColaboradores(limit, offset int) ([]entities.Colaborador, int64, error)
	GetColaboradorByID(id uuid.UUID) (*entities.Colaborador, error)
	CreateColaborador(colaborador *entities.Colaborador) error
	UpdateColaborador(id uuid.UUID, colaborador *entities.Colaborador) error
	DeleteColaborador(id uuid.UUID) error
	ListColaboradores(filters map[string]interface{}, limit, offset int) ([]entities.Colaborador, int64, error)
}

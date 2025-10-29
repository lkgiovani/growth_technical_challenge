package departamento

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/repository"
)

type usecase struct {
	departamentoRepo repository.DepartamentoRepository
	colaboradorRepo  repository.ColaboradorRepository
}

func NewUseCase(departamentoRepo repository.DepartamentoRepository, colaboradorRepo repository.ColaboradorRepository) UseCase {
	return &usecase{
		departamentoRepo: departamentoRepo,
		colaboradorRepo:  colaboradorRepo,
	}
}

func (u *usecase) GetAllDepartamentos(limit, offset int) ([]entities.Departamento, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return u.departamentoRepo.FindAll(limit, offset)
}

func (u *usecase) GetDepartamentoByID(id uuid.UUID) (*entities.Departamento, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID de departamento inválido")
	}
	return u.departamentoRepo.FindByID(id)
}

func (u *usecase) GetDepartamentoWithHierarchy(id uuid.UUID) (*entities.Departamento, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID de departamento inválido")
	}
	return u.departamentoRepo.FindByIDWithHierarchy(id)
}

func (u *usecase) CreateDepartamento(departamento *entities.Departamento) error {
	if err := u.validateDepartamento(departamento); err != nil {
		return err
	}

	gerente, err := u.colaboradorRepo.FindByID(departamento.GerenteID)
	if err != nil {
		return errors.New("gerente não encontrado")
	}
	if gerente == nil {
		return errors.New("gerente não existe")
	}

	if gerente.DepartamentoID != departamento.ID && departamento.ID != uuid.Nil {
		if gerente.DepartamentoID != departamento.ID {
			return errors.New("gerente deve pertencer ao mesmo departamento")
		}
	}

	if departamento.DepartamentoSuperiorID != nil && *departamento.DepartamentoSuperiorID != uuid.Nil {
		parentDept, err := u.departamentoRepo.FindByID(*departamento.DepartamentoSuperiorID)
		if err != nil {
			return errors.New("departamento superior não encontrado")
		}
		if parentDept == nil {
			return errors.New("departamento superior não existe")
		}

		if departamento.ID != uuid.Nil {
			hasCycle, err := u.departamentoRepo.HasCycle(departamento.ID, *departamento.DepartamentoSuperiorID)
			if err != nil {
				return err
			}
			if hasCycle {
				return errors.New("não é possível criar ciclo na hierarquia de departamentos")
			}
		}
	}

	departamento.Nome = strings.TrimSpace(departamento.Nome)
	return u.departamentoRepo.Create(departamento)
}

func (u *usecase) UpdateDepartamento(id uuid.UUID, departamento *entities.Departamento) error {
	if id == uuid.Nil {
		return errors.New("ID de departamento inválido")
	}

	existingDept, err := u.departamentoRepo.FindByID(id)
	if err != nil {
		return err
	}

	if errValidate := u.validateDepartamento(departamento); errValidate != nil {
		return errValidate
	}

	gerente, err := u.colaboradorRepo.FindByID(departamento.GerenteID)
	if err != nil {
		return errors.New("gerente não encontrado")
	}
	if gerente == nil {
		return errors.New("gerente não existe")
	}

	if gerente.DepartamentoID != id {
		return errors.New("gerente deve pertencer ao mesmo departamento")
	}

	if departamento.DepartamentoSuperiorID != nil && *departamento.DepartamentoSuperiorID != uuid.Nil {
		if *departamento.DepartamentoSuperiorID == id {
			return errors.New("departamento não pode ser seu próprio superior")
		}

		parentDept, err := u.departamentoRepo.FindByID(*departamento.DepartamentoSuperiorID)
		if err != nil {
			return errors.New("departamento superior não encontrado")
		}
		if parentDept == nil {
			return errors.New("departamento superior não existe")
		}

		hasCycle, err := u.departamentoRepo.HasCycle(id, *departamento.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if hasCycle {
			return errors.New("não é possível criar ciclo na hierarquia de departamentos")
		}
	}

	existingDept.Nome = strings.TrimSpace(departamento.Nome)
	existingDept.GerenteID = departamento.GerenteID
	existingDept.DepartamentoSuperiorID = departamento.DepartamentoSuperiorID

	return u.departamentoRepo.Update(existingDept)
}

func (u *usecase) DeleteDepartamento(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID de departamento inválido")
	}

	_, err := u.departamentoRepo.FindByID(id)
	if err != nil {
		return err
	}

	return u.departamentoRepo.Delete(id)
}

func (u *usecase) ListDepartamentos(filters map[string]interface{}, limit, offset int) ([]entities.Departamento, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return u.departamentoRepo.FindByFilters(filters, limit, offset)
}

func (u *usecase) GetGerenteColaboradores(gerenteID uuid.UUID) ([]entities.Colaborador, error) {
	if gerenteID == uuid.Nil {
		return nil, errors.New("ID de gerente inválido")
	}

	gerente, err := u.colaboradorRepo.FindByID(gerenteID)
	if err != nil {
		return nil, err
	}
	if gerente == nil {
		return nil, errors.New("gerente não encontrado")
	}

	departmentIDs, err := u.departamentoRepo.FindAllSubDepartmentIDs(gerente.DepartamentoID)
	if err != nil {
		return nil, err
	}

	return u.colaboradorRepo.FindByDepartmentIDs(departmentIDs)
}

func (u *usecase) validateDepartamento(departamento *entities.Departamento) error {
	if departamento == nil {
		return errors.New("departamento não pode ser nulo")
	}

	departamento.Nome = strings.TrimSpace(departamento.Nome)
	if departamento.Nome == "" {
		return errors.New("nome é obrigatório")
	}

	if len(departamento.Nome) < 3 {
		return errors.New("nome deve ter pelo menos 3 caracteres")
	}

	if len(departamento.Nome) > 255 {
		return errors.New("nome não pode exceder 255 caracteres")
	}

	if departamento.GerenteID == uuid.Nil {
		return errors.New("ID do gerente é obrigatório")
	}

	return nil
}

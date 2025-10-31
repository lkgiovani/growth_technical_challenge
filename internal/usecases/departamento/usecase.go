package departamento

import (
	"strings"

	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/repository"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"github.com/lkgiovani/growth_technical_challenge/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type usecase struct {
	departamentoRepo repository.DepartamentoRepository
	colaboradorRepo  repository.ColaboradorRepository
	db               *gorm.DB
	logger           logger.Logger
}

func NewUseCase(departamentoRepo repository.DepartamentoRepository, colaboradorRepo repository.ColaboradorRepository, db *gorm.DB, log logger.Logger) UseCase {
	return &usecase{
		departamentoRepo: departamentoRepo,
		colaboradorRepo:  colaboradorRepo,
		db:               db,
		logger:           log,
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
		return nil, app.Errorf(app.EINVALID, "invalid department ID")
	}
	return u.departamentoRepo.FindByID(id)
}

func (u *usecase) GetDepartamentoWithHierarchy(id uuid.UUID) (*entities.Departamento, error) {
	if id == uuid.Nil {
		return nil, app.Errorf(app.EINVALID, "invalid department ID")
	}
	return u.departamentoRepo.FindByIDWithHierarchy(id)
}

func (u *usecase) CreateDepartamento(departamento *entities.Departamento) error {
	departamento.ID = uuid.Nil

	if err := u.validateDepartamento(departamento); err != nil {
		u.logger.Warn("Department validation failed", zap.Error(err), zap.String("name", departamento.Nome))
		return err
	}

	u.logger.Debug("Creating department", zap.String("name", departamento.Nome), zap.String("managerID", departamento.GerenteID.String()))

	gerente, err := u.colaboradorRepo.FindByID(departamento.GerenteID)
	if err != nil {
		u.logger.Error("Failed to find manager", zap.Error(err), zap.String("managerID", departamento.GerenteID.String()))
		return err
	}
	if gerente == nil {
		u.logger.Warn("Manager not found", zap.String("managerID", departamento.GerenteID.String()))
		return app.Errorf(app.ENOTFOUND, "manager not found")
	}

	if departamento.DepartamentoSuperiorID != nil && *departamento.DepartamentoSuperiorID != uuid.Nil {
		parentDept, err := u.departamentoRepo.FindByID(*departamento.DepartamentoSuperiorID)
		if err != nil {
			u.logger.Error("Failed to find parent department", zap.Error(err), zap.String("parentID", departamento.DepartamentoSuperiorID.String()))
			return err
		}
		if parentDept == nil {
			u.logger.Warn("Parent department not found", zap.String("parentID", departamento.DepartamentoSuperiorID.String()))
			return app.Errorf(app.ENOTFOUND, "parent department not found")
		}
	}

	departamento.Nome = strings.TrimSpace(departamento.Nome)
	return u.departamentoRepo.Create(departamento)
}

func (u *usecase) CreateDepartamentoWithGerente(departamento *entities.Departamento, gerente *entities.Colaborador) error {
	if departamento == nil {
		u.logger.Warn("Attempt to create null department")
		return app.Errorf(app.EINVALID, "department cannot be null")
	}

	if gerente == nil {
		u.logger.Warn("Attempt to create department with null manager")
		return app.Errorf(app.EINVALID, "manager cannot be null")
	}

	u.logger.Info("Creating department with new manager", zap.String("departmentName", departamento.Nome), zap.String("managerName", gerente.Nome))

	departamento.ID = uuid.Nil
	gerente.ID = uuid.Nil

	departamento.Nome = strings.TrimSpace(departamento.Nome)
	if departamento.Nome == "" {
		return app.Errorf(app.EINVALID, "name is required")
	}

	if len(departamento.Nome) < 3 {
		return app.Errorf(app.EINVALID, "name must be at least 3 characters")
	}

	if len(departamento.Nome) > 255 {
		return app.Errorf(app.EINVALID, "name cannot exceed 255 characters")
	}

	if err := u.validateColaborador(gerente); err != nil {
		return err
	}

	gerente.CPF = utils.NormalizeCPF(gerente.CPF)
	if !utils.ValidateCPF(gerente.CPF) {
		return app.Errorf(app.EINVALID, "invalid CPF")
	}

	existingColaborador, err := u.colaboradorRepo.FindByCPF(gerente.CPF)
	if err != nil {
		return err
	}
	if existingColaborador != nil {
		return app.Errorf(app.EDUPLICATION, "CPF already exists")
	}

	if gerente.RG != nil && *gerente.RG != "" {
		normalizedRG := utils.NormalizeRG(*gerente.RG)
		gerente.RG = &normalizedRG

		existingRG, errRG := u.colaboradorRepo.FindByRG(*gerente.RG)
		if errRG != nil {
			return errRG
		}
		if existingRG != nil {
			return app.Errorf(app.EDUPLICATION, "RG already exists")
		}
	}

	if departamento.DepartamentoSuperiorID != nil && *departamento.DepartamentoSuperiorID != uuid.Nil {
		parentDept, err := u.departamentoRepo.FindByID(*departamento.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if parentDept == nil {
			return app.Errorf(app.ENOTFOUND, "parent department not found")
		}
	}

	gerente.Nome = strings.TrimSpace(gerente.Nome)

	return u.db.Transaction(func(tx *gorm.DB) error {
		txColabRepo := repository.NewColaboradorRepository(tx, nil, u.logger)
		txDeptRepo := repository.NewDepartamentoRepository(tx, nil, u.logger)

		if err := txColabRepo.Create(gerente); err != nil {
			return err
		}

		departamento.GerenteID = gerente.ID
		gerente.DepartamentoID = departamento.ID

		if err := txDeptRepo.Create(departamento); err != nil {
			return err
		}

		gerente.DepartamentoID = departamento.ID
		if err := txColabRepo.Update(gerente); err != nil {
			return err
		}

		return nil
	})
}

func (u *usecase) validateColaborador(colaborador *entities.Colaborador) error {
	if colaborador == nil {
		return app.Errorf(app.EINVALID, "employee cannot be null")
	}

	colaborador.Nome = strings.TrimSpace(colaborador.Nome)
	if colaborador.Nome == "" {
		return app.Errorf(app.EINVALID, "manager name is required")
	}

	if len(colaborador.Nome) < 3 {
		return app.Errorf(app.EINVALID, "manager name must be at least 3 characters")
	}

	if len(colaborador.Nome) > 255 {
		return app.Errorf(app.EINVALID, "manager name cannot exceed 255 characters")
	}

	colaborador.CPF = utils.NormalizeCPF(colaborador.CPF)
	if colaborador.CPF == "" {
		return app.Errorf(app.EINVALID, "manager CPF is required")
	}

	if !utils.ValidateCPF(colaborador.CPF) {
		return app.Errorf(app.EINVALID, "invalid CPF")
	}

	return nil
}

func (u *usecase) UpdateDepartamento(id uuid.UUID, departamento *entities.Departamento) error {
	if id == uuid.Nil {
		u.logger.Warn("Invalid department ID for update", zap.String("id", id.String()))
		return app.Errorf(app.EINVALID, "invalid department ID")
	}

	u.logger.Debug("Updating department", zap.String("id", id.String()))

	existingDept, err := u.departamentoRepo.FindByID(id)
	if err != nil {
		u.logger.Error("Failed to find department for update", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	if departamento.Nome != "" {
		nome := strings.TrimSpace(departamento.Nome)
		if nome == "" {
			return app.Errorf(app.EINVALID, "department name cannot be empty")
		}
		if len(nome) < 3 {
			return app.Errorf(app.EINVALID, "department name must be at least 3 characters")
		}
		if len(nome) > 255 {
			return app.Errorf(app.EINVALID, "department name cannot exceed 255 characters")
		}
		existingDept.Nome = nome
	}

	if departamento.GerenteID != uuid.Nil {
		gerente, err := u.colaboradorRepo.FindByID(departamento.GerenteID)
		if err != nil {
			return err
		}
		if gerente == nil {
			return app.Errorf(app.ENOTFOUND, "manager not found")
		}

		if gerente.DepartamentoID != id {
			return app.Errorf(app.EINVALID, "manager must belong to the same department")
		}

		existingDept.GerenteID = departamento.GerenteID
	}

	if departamento.DepartamentoSuperiorID != nil {
		if *departamento.DepartamentoSuperiorID == uuid.Nil {
			existingDept.DepartamentoSuperiorID = nil
		} else {
			if *departamento.DepartamentoSuperiorID == id {
				return app.Errorf(app.EINVALID, "department cannot be its own parent")
			}

			parentDept, err := u.departamentoRepo.FindByID(*departamento.DepartamentoSuperiorID)
			if err != nil {
				return err
			}
			if parentDept == nil {
				return app.Errorf(app.ENOTFOUND, "parent department not found")
			}

			hasCycle, err := u.departamentoRepo.HasCycle(id, *departamento.DepartamentoSuperiorID)
			if err != nil {
				return err
			}
			if hasCycle {
				return app.Errorf(app.ECONFLICT, "cannot create cycle in department hierarchy")
			}

			existingDept.DepartamentoSuperiorID = departamento.DepartamentoSuperiorID
		}
	}

	return u.departamentoRepo.Update(existingDept)
}

func (u *usecase) DeleteDepartamento(id uuid.UUID) error {
	if id == uuid.Nil {
		u.logger.Warn("Invalid department ID for deletion", zap.String("id", id.String()))
		return app.Errorf(app.EINVALID, "invalid department ID")
	}

	u.logger.Debug("Deleting department", zap.String("id", id.String()))

	departamento, err := u.departamentoRepo.FindByID(id)
	if err != nil {
		u.logger.Error("Failed to find department for deletion", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	subdepartamentos, err := u.departamentoRepo.FindSubDepartments(id)
	if err != nil {
		u.logger.Error("Failed to find subdepartments", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	if len(subdepartamentos) > 0 {
		u.logger.Info("Department has subdepartments, reassigning parent", zap.String("id", id.String()), zap.Int("subdepartmentCount", len(subdepartamentos)))
		if departamento.DepartamentoSuperiorID == nil {
			u.logger.Warn("Cannot delete department with subdepartments without parent", zap.String("id", id.String()))
			return app.Errorf(app.EINVALID, "cannot delete department with subdepartments without a parent department")
		}

		for _, sub := range subdepartamentos {
			sub.DepartamentoSuperiorID = departamento.DepartamentoSuperiorID
			if err := u.departamentoRepo.Update(&sub); err != nil {
				u.logger.Error("Failed to reassign subdepartment parent", zap.Error(err), zap.String("subdepartmentID", sub.ID.String()))
				return err
			}
		}
	}

	u.logger.Info("Department deleted successfully", zap.String("id", id.String()))
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
		u.logger.Warn("Invalid manager ID", zap.String("id", gerenteID.String()))
		return nil, app.Errorf(app.EINVALID, "invalid manager ID")
	}

	u.logger.Debug("Fetching colaboradores for manager", zap.String("managerID", gerenteID.String()))

	gerente, err := u.colaboradorRepo.FindByID(gerenteID)
	if err != nil {
		u.logger.Error("Failed to find manager", zap.Error(err), zap.String("managerID", gerenteID.String()))
		return nil, err
	}
	if gerente == nil {
		u.logger.Warn("Manager not found", zap.String("managerID", gerenteID.String()))
		return nil, app.Errorf(app.ENOTFOUND, "gerente não encontrado")
	}

	departmentIDs, err := u.departamentoRepo.FindAllSubDepartmentIDs(gerente.DepartamentoID)
	if err != nil {
		u.logger.Error("Failed to find subdepartment IDs", zap.Error(err), zap.String("departmentID", gerente.DepartamentoID.String()))
		return nil, err
	}

	u.logger.Debug("Found department hierarchy", zap.String("managerID", gerenteID.String()), zap.Int("departmentCount", len(departmentIDs)))

	return u.colaboradorRepo.FindByDepartmentIDs(departmentIDs)
}

func (u *usecase) validateDepartamento(departamento *entities.Departamento) error {
	if departamento == nil {
		return app.Errorf(app.EINVALID, "department cannot be null")
	}

	departamento.Nome = strings.TrimSpace(departamento.Nome)
	if departamento.Nome == "" {
		return app.Errorf(app.EINVALID, "name is required")
	}

	if len(departamento.Nome) < 3 {
		return app.Errorf(app.EINVALID, "name must be at least 3 characters")
	}

	if len(departamento.Nome) > 255 {
		return app.Errorf(app.EINVALID, "name cannot exceed 255 characters")
	}

	if departamento.GerenteID == uuid.Nil {
		return app.Errorf(app.EINVALID, "manager ID is required")
	}

	return nil
}

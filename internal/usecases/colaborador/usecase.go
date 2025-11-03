package colaborador

import (
	"strings"

	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/repository"
	"github.com/lkgiovani/growth_technical_challenge/internal/services/cache_service"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"github.com/lkgiovani/growth_technical_challenge/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type usecase struct {
	colaboradorRepo  repository.ColaboradorRepository
	departamentoRepo repository.DepartamentoRepository
	cacheService     cache_service.CacheService
	db               *gorm.DB
	logger           logger.Logger
}

func NewUseCase(colaboradorRepo repository.ColaboradorRepository, departamentoRepo repository.DepartamentoRepository, cacheSvc cache_service.CacheService, db *gorm.DB, log logger.Logger) UseCase {
	return &usecase{
		colaboradorRepo:  colaboradorRepo,
		departamentoRepo: departamentoRepo,
		cacheService:     cacheSvc,
		db:               db,
		logger:           log,
	}
}

func (u *usecase) GetAllColaboradores(limit, offset int) ([]entities.Colaborador, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return u.colaboradorRepo.FindAll(limit, offset)
}

func (u *usecase) GetColaboradorByID(id uuid.UUID) (*entities.Colaborador, error) {
	if id == uuid.Nil {
		return nil, app.Errorf(app.EINVALID, "ID de colaborador inválido")
	}
	return u.colaboradorRepo.FindByID(id)
}

func (u *usecase) CreateColaborador(colaborador *entities.Colaborador) error {
	colaborador.ID = uuid.Nil

	if err := u.validateColaborador(colaborador); err != nil {
		u.logger.Warn("Colaborador validation failed", zap.Error(err), zap.String("name", colaborador.Nome))
		return err
	}

	colaborador.CPF = utils.NormalizeCPF(colaborador.CPF)
	if !utils.ValidateCPF(colaborador.CPF) {
		u.logger.Warn("Invalid CPF", zap.String("cpf", colaborador.CPF))
		return app.Errorf(app.EINVALID, "CPF inválido")
	}

	existingColaborador, err := u.colaboradorRepo.FindByCPF(colaborador.CPF)
	if err != nil {
		u.logger.Error("Failed to check existing CPF", zap.Error(err), zap.String("cpf", colaborador.CPF))
		return err
	}
	if existingColaborador != nil {
		u.logger.Warn("CPF already exists", zap.String("cpf", colaborador.CPF))
		return app.Errorf(app.EDUPLICATION, "CPF já existe")
	}

	if colaborador.RG != nil && *colaborador.RG != "" {
		normalizedRG := utils.NormalizeRG(*colaborador.RG)
		colaborador.RG = &normalizedRG

		existingRG, errRG := u.colaboradorRepo.FindByRG(*colaborador.RG)
		if errRG != nil {
			u.logger.Error("Failed to check existing RG", zap.Error(errRG), zap.String("rg", *colaborador.RG))
			return errRG
		}
		if existingRG != nil {
			u.logger.Warn("RG already exists", zap.String("rg", *colaborador.RG))
			return app.Errorf(app.EDUPLICATION, "RG já existe")
		}
	}

	departamento, errDept := u.departamentoRepo.FindByID(colaborador.DepartamentoID)
	if errDept != nil {
		u.logger.Error("Failed to find department", zap.Error(errDept), zap.String("departmentID", colaborador.DepartamentoID.String()))
		return errDept
	}
	if departamento == nil {
		u.logger.Warn("Department not found", zap.String("departmentID", colaborador.DepartamentoID.String()))
		return app.Errorf(app.ENOTFOUND, "departamento não encontrado")
	}

	colaborador.Nome = strings.TrimSpace(colaborador.Nome)

	u.logger.Debug("Creating colaborador in repository", zap.String("name", colaborador.Nome), zap.String("cpf", colaborador.CPF))

	if err := u.colaboradorRepo.Create(colaborador); err != nil {
		return err
	}

	if u.cacheService != nil {
		if err := u.cacheService.InvalidateColaboradoresByDepartment(colaborador.DepartamentoID); err != nil {
			u.logger.Warn("Failed to invalidate colaboradores cache", zap.Error(err))
		}
	}

	return nil
}

func (u *usecase) CreateColaboradorWithDepartamento(colaborador *entities.Colaborador, departamento *entities.Departamento) error {
	if colaborador == nil {
		u.logger.Warn("Attempt to create null colaborador")
		return app.Errorf(app.EINVALID, "colaborador cannot be null")
	}

	if departamento == nil {
		u.logger.Warn("Attempt to create colaborador with null departamento")
		return app.Errorf(app.EINVALID, "departamento cannot be null")
	}

	u.logger.Info("Creating colaborador with new departamento", zap.String("colaboradorName", colaborador.Nome), zap.String("departamentoName", departamento.Nome))

	colaborador.ID = uuid.Nil
	departamento.ID = uuid.Nil

	colaborador.Nome = strings.TrimSpace(colaborador.Nome)
	if colaborador.Nome == "" {
		return app.Errorf(app.EINVALID, "nome é obrigatório")
	}

	if len(colaborador.Nome) < 3 {
		return app.Errorf(app.EINVALID, "nome deve ter pelo menos 3 caracteres")
	}

	if len(colaborador.Nome) > 255 {
		return app.Errorf(app.EINVALID, "nome não pode exceder 255 caracteres")
	}

	departamento.Nome = strings.TrimSpace(departamento.Nome)
	if departamento.Nome == "" {
		return app.Errorf(app.EINVALID, "nome do departamento é obrigatório")
	}

	if len(departamento.Nome) < 3 {
		return app.Errorf(app.EINVALID, "nome do departamento deve ter pelo menos 3 caracteres")
	}

	if len(departamento.Nome) > 255 {
		return app.Errorf(app.EINVALID, "nome do departamento não pode exceder 255 caracteres")
	}

	colaborador.CPF = utils.NormalizeCPF(colaborador.CPF)
	if !utils.ValidateCPF(colaborador.CPF) {
		return app.Errorf(app.EINVALID, "CPF inválido")
	}

	existingColaborador, err := u.colaboradorRepo.FindByCPF(colaborador.CPF)
	if err != nil {
		return err
	}
	if existingColaborador != nil {
		return app.Errorf(app.EDUPLICATION, "CPF já existe")
	}

	if colaborador.RG != nil && *colaborador.RG != "" {
		normalizedRG := utils.NormalizeRG(*colaborador.RG)
		colaborador.RG = &normalizedRG

		existingRG, errRG := u.colaboradorRepo.FindByRG(*colaborador.RG)
		if errRG != nil {
			return errRG
		}
		if existingRG != nil {
			return app.Errorf(app.EDUPLICATION, "RG já existe")
		}
	}

	gerente, err := u.colaboradorRepo.FindByID(departamento.GerenteID)
	if err != nil {
		u.logger.Error("Failed to find gerente", zap.Error(err), zap.String("gerenteID", departamento.GerenteID.String()))
		return err
	}
	if gerente == nil {
		u.logger.Warn("Gerente not found", zap.String("gerenteID", departamento.GerenteID.String()))
		return app.Errorf(app.ENOTFOUND, "gerente não encontrado")
	}

	if departamento.DepartamentoSuperiorID != nil && *departamento.DepartamentoSuperiorID != uuid.Nil {
		parentDept, err := u.departamentoRepo.FindByID(*departamento.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if parentDept == nil {
			return app.Errorf(app.ENOTFOUND, "departamento superior não encontrado")
		}
	}

	return u.db.Transaction(func(tx *gorm.DB) error {
		txDeptRepo := repository.NewDepartamentoRepository(tx, u.logger)
		txColabRepo := repository.NewColaboradorRepository(tx, u.logger)

		if err := txDeptRepo.Create(departamento); err != nil {
			return err
		}

		colaborador.DepartamentoID = departamento.ID

		if err := txColabRepo.Create(colaborador); err != nil {
			return err
		}

		if u.cacheService != nil {
			if err := u.cacheService.InvalidateColaboradoresByDepartment(departamento.ID); err != nil {
				u.logger.Warn("Failed to invalidate colaboradores cache", zap.Error(err))
			}
			if err := u.cacheService.InvalidateAllDepartmentCache(departamento.ID, departamento.DepartamentoSuperiorID); err != nil {
				u.logger.Warn("Failed to invalidate department cache", zap.Error(err))
			}
		}

		return nil
	})
}

func (u *usecase) UpdateColaborador(id uuid.UUID, colaborador *entities.Colaborador) error {
	if id == uuid.Nil {
		u.logger.Warn("Invalid colaborador ID", zap.String("id", id.String()))
		return app.Errorf(app.EINVALID, "ID de colaborador inválido")
	}

	u.logger.Debug("Updating colaborador", zap.String("id", id.String()))

	existingColaborador, err := u.colaboradorRepo.FindByID(id)
	if err != nil {
		u.logger.Error("Failed to find colaborador", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	if colaborador.Nome != "" {
		nome := strings.TrimSpace(colaborador.Nome)
		if nome == "" {
			return app.Errorf(app.EINVALID, "nome não pode ser vazio")
		}
		if len(nome) < 3 {
			return app.Errorf(app.EINVALID, "nome deve ter pelo menos 3 caracteres")
		}
		if len(nome) > 255 {
			return app.Errorf(app.EINVALID, "nome não pode exceder 255 caracteres")
		}
		existingColaborador.Nome = nome
	}

	if colaborador.CPF != "" {
		cpfNormalizado := utils.NormalizeCPF(colaborador.CPF)
		if !utils.ValidateCPF(cpfNormalizado) {
			u.logger.Warn("Invalid CPF on update", zap.String("cpf", cpfNormalizado), zap.String("id", id.String()))
			return app.Errorf(app.EINVALID, "CPF inválido")
		}

		cpfColaborador, err := u.colaboradorRepo.FindByCPF(cpfNormalizado)
		if err != nil {
			u.logger.Error("Failed to check existing CPF on update", zap.Error(err), zap.String("cpf", cpfNormalizado))
			return err
		}
		if cpfColaborador != nil && cpfColaborador.ID != id {
			u.logger.Warn("CPF already exists on update", zap.String("cpf", cpfNormalizado), zap.String("id", id.String()))
			return app.Errorf(app.EDUPLICATION, "CPF já existe")
		}

		existingColaborador.CPF = cpfNormalizado
	}

	if colaborador.RG != nil && *colaborador.RG != "" {
		normalizedRG := utils.NormalizeRG(*colaborador.RG)
		colaborador.RG = &normalizedRG

		rgColaborador, errRG := u.colaboradorRepo.FindByRG(*colaborador.RG)
		if errRG != nil {
			return errRG
		}
		if rgColaborador != nil && rgColaborador.ID != id {
			return app.Errorf(app.EDUPLICATION, "RG já existe")
		}

		existingColaborador.RG = colaborador.RG
	}

	if colaborador.DepartamentoID != uuid.Nil {
		departamento, errDept := u.departamentoRepo.FindByID(colaborador.DepartamentoID)
		if errDept != nil {
			return errDept
		}
		if departamento == nil {
			return app.Errorf(app.ENOTFOUND, "departamento não encontrado")
		}

		existingColaborador.DepartamentoID = colaborador.DepartamentoID
	}

	oldDepartmentID := existingColaborador.DepartamentoID
	newDepartmentID := colaborador.DepartamentoID

	if err := u.colaboradorRepo.Update(existingColaborador); err != nil {
		return err
	}

	if u.cacheService != nil {
		if err := u.cacheService.InvalidateColaboradoresByDepartment(oldDepartmentID); err != nil {
			u.logger.Warn("Failed to invalidate old department cache", zap.Error(err))
		}

		if newDepartmentID != uuid.Nil && newDepartmentID != oldDepartmentID {
			if err := u.cacheService.InvalidateColaboradoresByDepartment(newDepartmentID); err != nil {
				u.logger.Warn("Failed to invalidate new department cache", zap.Error(err))
			}
		}
	}

	return nil
}

func (u *usecase) DeleteColaborador(id uuid.UUID) error {
	if id == uuid.Nil {
		u.logger.Warn("Invalid colaborador ID for deletion", zap.String("id", id.String()))
		return app.Errorf(app.EINVALID, "ID de colaborador inválido")
	}

	u.logger.Debug("Deleting colaborador", zap.String("id", id.String()))

	colaborador, err := u.colaboradorRepo.FindByID(id)
	if err != nil {
		u.logger.Error("Failed to find colaborador for deletion", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	if err := u.colaboradorRepo.Delete(id); err != nil {
		return err
	}

	if u.cacheService != nil {
		if err := u.cacheService.InvalidateColaboradoresByDepartment(colaborador.DepartamentoID); err != nil {
			u.logger.Warn("Failed to invalidate colaboradores cache", zap.Error(err))
		}
	}

	return nil
}

func (u *usecase) ListColaboradores(filters map[string]interface{}, limit, offset int) ([]entities.Colaborador, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return u.colaboradorRepo.FindByFilters(filters, limit, offset)
}

func (u *usecase) validateColaborador(colaborador *entities.Colaborador) error {
	if colaborador == nil {
		return app.Errorf(app.EINVALID, "colaborador não pode ser nulo")
	}

	colaborador.Nome = strings.TrimSpace(colaborador.Nome)
	if colaborador.Nome == "" {
		return app.Errorf(app.EINVALID, "nome é obrigatório")
	}

	if len(colaborador.Nome) < 3 {
		return app.Errorf(app.EINVALID, "nome deve ter pelo menos 3 caracteres")
	}

	if len(colaborador.Nome) > 255 {
		return app.Errorf(app.EINVALID, "nome não pode exceder 255 caracteres")
	}

	colaborador.CPF = strings.TrimSpace(colaborador.CPF)
	if colaborador.CPF == "" {
		return app.Errorf(app.EINVALID, "CPF é obrigatório")
	}

	if colaborador.DepartamentoID == uuid.Nil {
		return app.Errorf(app.EINVALID, "ID do departamento é obrigatório")
	}

	return nil
}

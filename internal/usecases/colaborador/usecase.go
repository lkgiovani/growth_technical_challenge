package colaborador

import (
	"strings"

	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/internal/repository"
	"github.com/lkgiovani/growth_technical_challenge/pkg/utils"
)

type usecase struct {
	colaboradorRepo  repository.ColaboradorRepository
	departamentoRepo repository.DepartamentoRepository
}

func NewUseCase(colaboradorRepo repository.ColaboradorRepository, departamentoRepo repository.DepartamentoRepository) UseCase {
	return &usecase{
		colaboradorRepo:  colaboradorRepo,
		departamentoRepo: departamentoRepo,
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
	if err := u.validateColaborador(colaborador); err != nil {
		return err
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

	departamento, errDept := u.departamentoRepo.FindByID(colaborador.DepartamentoID)
	if errDept != nil {
		return errDept
	}
	if departamento == nil {
		return app.Errorf(app.ENOTFOUND, "departamento não encontrado")
	}

	colaborador.Nome = strings.TrimSpace(colaborador.Nome)

	return u.colaboradorRepo.Create(colaborador)
}

func (u *usecase) UpdateColaborador(id uuid.UUID, colaborador *entities.Colaborador) error {
	if id == uuid.Nil {
		return app.Errorf(app.EINVALID, "ID de colaborador inválido")
	}

	existingColaborador, err := u.colaboradorRepo.FindByID(id)
	if err != nil {
		return err
	}

	if errValidate := u.validateColaborador(colaborador); errValidate != nil {
		return errValidate
	}

	colaborador.CPF = utils.NormalizeCPF(colaborador.CPF)
	if !utils.ValidateCPF(colaborador.CPF) {
		return app.Errorf(app.EINVALID, "CPF inválido")
	}

	cpfColaborador, err := u.colaboradorRepo.FindByCPF(colaborador.CPF)
	if err != nil {
		return err
	}
	if cpfColaborador != nil && cpfColaborador.ID != id {
		return app.Errorf(app.EDUPLICATION, "CPF já existe")
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
	}

	departamento, errDept := u.departamentoRepo.FindByID(colaborador.DepartamentoID)
	if errDept != nil {
		return errDept
	}
	if departamento == nil {
		return app.Errorf(app.ENOTFOUND, "departamento não encontrado")
	}

	existingColaborador.Nome = strings.TrimSpace(colaborador.Nome)
	existingColaborador.CPF = colaborador.CPF
	existingColaborador.RG = colaborador.RG
	existingColaborador.DepartamentoID = colaborador.DepartamentoID

	return u.colaboradorRepo.Update(existingColaborador)
}

func (u *usecase) DeleteColaborador(id uuid.UUID) error {
	if id == uuid.Nil {
		return app.Errorf(app.EINVALID, "ID de colaborador inválido")
	}

	_, err := u.colaboradorRepo.FindByID(id)
	if err != nil {
		return err
	}

	return u.colaboradorRepo.Delete(id)
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

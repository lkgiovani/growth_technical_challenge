package repository

import (
	"errors"

	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"gorm.io/gorm"
)

type colaboradorRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewColaboradorRepository(db *gorm.DB, log logger.Logger) ColaboradorRepository {
	return &colaboradorRepository{
		db:     db,
		logger: log,
	}
}

func (r *colaboradorRepository) FindAll(limit, offset int) ([]entities.Colaborador, int64, error) {
	var colaboradores []entities.Colaborador
	var total int64

	if err := r.db.Model(&entities.Colaborador{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Departamento").Limit(limit).Offset(offset).Find(&colaboradores)
	return colaboradores, total, result.Error
}

func (r *colaboradorRepository) FindByID(id uuid.UUID) (*entities.Colaborador, error) {
	var colaborador entities.Colaborador
	result := r.db.Preload("Departamento.Gerente").First(&colaborador, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app.Errorf(app.ENOTFOUND, "employee not found")
		}
		return nil, result.Error
	}
	return &colaborador, nil
}

func (r *colaboradorRepository) FindByCPF(cpf string) (*entities.Colaborador, error) {
	var colaborador entities.Colaborador
	result := r.db.Where("cpf = ?", cpf).First(&colaborador)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &colaborador, nil
}

func (r *colaboradorRepository) FindByRG(rg string) (*entities.Colaborador, error) {
	var colaborador entities.Colaborador
	result := r.db.Where("rg = ?", rg).First(&colaborador)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &colaborador, nil
}

func (r *colaboradorRepository) FindByFilters(filters map[string]interface{}, limit, offset int) ([]entities.Colaborador, int64, error) {
	var colaboradores []entities.Colaborador
	var total int64

	query := r.db.Model(&entities.Colaborador{})

	if nome, ok := filters["nome"].(string); ok && nome != "" {
		query = query.Where("nome ILIKE ?", "%"+nome+"%")
	}
	if cpf, ok := filters["cpf"].(string); ok && cpf != "" {
		query = query.Where("cpf = ?", cpf)
	}
	if rg, ok := filters["rg"].(string); ok && rg != "" {
		query = query.Where("rg = ?", rg)
	}
	if deptID, ok := filters["departamento_id"].(string); ok && deptID != "" {
		if id, err := uuid.Parse(deptID); err == nil {
			query = query.Where("departamento_id = ?", id)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := query.Preload("Departamento").Limit(limit).Offset(offset).Find(&colaboradores)
	return colaboradores, total, result.Error
}

func (r *colaboradorRepository) FindByDepartmentIDs(departmentIDs []uuid.UUID) ([]entities.Colaborador, error) {
	var colaboradores []entities.Colaborador
	result := r.db.Preload("Departamento").Where("departamento_id IN ?", departmentIDs).Find(&colaboradores)
	return colaboradores, result.Error
}

func (r *colaboradorRepository) Create(colaborador *entities.Colaborador) error {
	return r.db.Create(colaborador).Error
}

func (r *colaboradorRepository) Update(colaborador *entities.Colaborador) error {
	return r.db.Model(colaborador).Updates(colaborador).Error
}

func (r *colaboradorRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&entities.Colaborador{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return app.Errorf(app.ENOTFOUND, "employee not found")
	}
	return nil
}

func (r *colaboradorRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.Colaborador{}).Count(&count)
	return count, result.Error
}

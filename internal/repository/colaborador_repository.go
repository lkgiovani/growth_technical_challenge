package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/pkg/cache"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

type colaboradorRepository struct {
	db     *gorm.DB
	cache  cache.Cache
	logger logger.Logger
}

func NewColaboradorRepository(db *gorm.DB, cacheClient cache.Cache, log logger.Logger) ColaboradorRepository {
	return &colaboradorRepository{
		db:     db,
		cache:  cacheClient,
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
	if r.cache != nil {
		ctx := context.Background()

		if len(departmentIDs) == 0 {
			return []entities.Colaborador{}, nil
		}

		if len(departmentIDs) == 1 {
			cacheKey := fmt.Sprintf("colaborador:dept:%s", departmentIDs[0].String())

			var colaboradores []entities.Colaborador
			err := cache.GetJSON(ctx, r.cache, cacheKey, &colaboradores)
			if err == nil {
				r.logger.Debug("Cache hit for colaboradores by department", zap.String("departmentID", departmentIDs[0].String()))
				return colaboradores, nil
			}

			if err != redis.Nil {
				r.logger.Warn("Cache error", zap.Error(err), zap.String("key", cacheKey))
			}

			result := r.db.Preload("Departamento").Where("departamento_id IN ?", departmentIDs).Find(&colaboradores)
			if result.Error != nil {
				return nil, result.Error
			}

			if err := cache.SetJSON(ctx, r.cache, cacheKey, colaboradores); err != nil {
				r.logger.Warn("Failed to cache colaboradores by department", zap.Error(err), zap.String("key", cacheKey))
			}

			return colaboradores, nil
		}

		result := make([]entities.Colaborador, 0)
		uncachedDeptIDs := make([]uuid.UUID, 0)

		for _, deptID := range departmentIDs {
			cacheKey := fmt.Sprintf("colaborador:dept:%s", deptID.String())

			var colaboradores []entities.Colaborador
			err := cache.GetJSON(ctx, r.cache, cacheKey, &colaboradores)
			if err == nil {
				result = append(result, colaboradores...)
			} else {
				if err != redis.Nil {
					r.logger.Warn("Cache error for department", zap.Error(err), zap.String("departmentID", deptID.String()))
				}
				uncachedDeptIDs = append(uncachedDeptIDs, deptID)
			}
		}

		if len(uncachedDeptIDs) > 0 {
			var cols []entities.Colaborador
			dbResult := r.db.Preload("Departamento").Where("departamento_id IN ?", uncachedDeptIDs).Find(&cols)
			if dbResult.Error != nil {
				return nil, dbResult.Error
			}

			deptColsMap := make(map[uuid.UUID][]entities.Colaborador)
			for _, col := range cols {
				deptColsMap[col.DepartamentoID] = append(deptColsMap[col.DepartamentoID], col)
			}

			for deptID, deptCols := range deptColsMap {
				cacheKey := fmt.Sprintf("colaborador:dept:%s", deptID.String())
				if err := cache.SetJSON(ctx, r.cache, cacheKey, deptCols); err != nil {
					r.logger.Warn("Failed to cache colaboradores for department", zap.Error(err), zap.String("departmentID", deptID.String()))
				}
			}

			result = append(result, cols...)
		}

		return result, nil
	}

	var colaboradores []entities.Colaborador
	result := r.db.Preload("Departamento").Where("departamento_id IN ?", departmentIDs).Find(&colaboradores)
	return colaboradores, result.Error
}

func (r *colaboradorRepository) Create(colaborador *entities.Colaborador) error {
	err := r.db.Create(colaborador).Error
	if err != nil {
		return err
	}

	if r.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("colaborador:dept:%s", colaborador.DepartamentoID.String())
		if err := r.cache.Del(ctx, cacheKey); err != nil {
			r.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("key", cacheKey))
		}
	}
	return nil
}

func (r *colaboradorRepository) Update(colaborador *entities.Colaborador) error {
	err := r.db.Model(colaborador).Updates(colaborador).Error
	if err != nil {
		return err
	}

	if r.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("colaborador:dept:%s", colaborador.DepartamentoID.String())
		if err := r.cache.Del(ctx, cacheKey); err != nil {
			r.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("key", cacheKey))
		}
	}
	return nil
}

func (r *colaboradorRepository) Delete(id uuid.UUID) error {
	col, err := r.FindByID(id)
	if err != nil {
		return err
	}

	result := r.db.Delete(&entities.Colaborador{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return app.Errorf(app.ENOTFOUND, "employee not found")
	}

	if r.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("colaborador:dept:%s", col.DepartamentoID.String())
		if err := r.cache.Del(ctx, cacheKey); err != nil {
			r.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("key", cacheKey))
		}
	}

	return nil
}

func (r *colaboradorRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.Colaborador{}).Count(&count)
	return count, result.Error
}

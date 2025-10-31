package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	app "github.com/lkgiovani/growth_technical_challenge"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/pkg/cache"
	"github.com/redis/go-redis/v9"
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
	db    *gorm.DB
	cache cache.Cache
}

func NewColaboradorRepository(db *gorm.DB, cacheClient cache.Cache) ColaboradorRepository {
	return &colaboradorRepository{
		db:    db,
		cache: cacheClient,
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
				return colaboradores, nil
			}

			if err != redis.Nil {
				fmt.Printf("Erro no cache: %v\n", err)
			}

			result := r.db.Preload("Departamento").Where("departamento_id IN ?", departmentIDs).Find(&colaboradores)
			if result.Error != nil {
				return nil, result.Error
			}

			if err := cache.SetJSON(ctx, r.cache, cacheKey, colaboradores); err != nil {
				fmt.Printf("Falha ao cachear colaboradores por departamento: %v\n", err)
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
					fmt.Printf("Erro no cache para dept %s: %v\n", deptID, err)
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
					fmt.Printf("Falha ao cachear colaboradores do dept %s: %v\n", deptID, err)
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
			fmt.Printf("Falha ao invalidar cache para %s: %v\n", cacheKey, err)
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
			fmt.Printf("Falha ao invalidar cache para %s: %v\n", cacheKey, err)
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
			fmt.Printf("Falha ao invalidar cache para %s: %v\n", cacheKey, err)
		}
	}

	return nil
}

func (r *colaboradorRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.Colaborador{}).Count(&count)
	return count, result.Error
}

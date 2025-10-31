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

type departamentoRepository struct {
	db     *gorm.DB
	cache  cache.Cache
	logger logger.Logger
}

func NewDepartamentoRepository(db *gorm.DB, cacheClient cache.Cache, log logger.Logger) DepartamentoRepository {
	return &departamentoRepository{
		db:     db,
		cache:  cacheClient,
		logger: log,
	}
}

func (r *departamentoRepository) FindAll(limit, offset int) ([]entities.Departamento, int64, error) {
	var departamentos []entities.Departamento
	var total int64

	if err := r.db.Model(&entities.Departamento{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Gerente").Preload("DepartamentoSuperior").
		Limit(limit).Offset(offset).Find(&departamentos)
	return departamentos, total, result.Error
}

func (r *departamentoRepository) FindByID(id uuid.UUID) (*entities.Departamento, error) {
	var departamento entities.Departamento
	result := r.db.Preload("Gerente").Preload("DepartamentoSuperior").
		First(&departamento, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app.Errorf(app.ENOTFOUND, "department not found")
		}
		return nil, result.Error
	}
	return &departamento, nil
}

func (r *departamentoRepository) FindByIDWithHierarchy(id uuid.UUID) (*entities.Departamento, error) {
	if r.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("departamento:hierarquia:%s", id.String())

		var departamento entities.Departamento
		err := cache.GetJSON(ctx, r.cache, cacheKey, &departamento)
		if err == nil {
			return &departamento, nil
		}

		if err != redis.Nil {
			r.logger.Warn("Cache error", zap.Error(err), zap.String("key", cacheKey))
		}
	}

	type DeptHierarchy struct {
		ID                     uuid.UUID  `json:"id"`
		Nome                   string     `json:"nome"`
		GerenteID              uuid.UUID  `json:"gerente_id"`
		GerenteNome            string     `json:"gerente_nome"`
		DepartamentoSuperiorID *uuid.UUID `json:"departamento_superior_id"`
		Level                  int        `json:"level"`
	}

	query := `
		WITH RECURSIVE department_tree AS (
			SELECT 
				d.id,
				d.nome,
				d.gerente_id,
				c.nome as gerente_nome,
				d.departamento_superior_id,
				0 as level
			FROM departamentos d
			LEFT JOIN colaboradores c ON d.gerente_id = c.id
			WHERE d.id = ? AND d.deleted_at IS NULL
			
			UNION ALL
			
			SELECT 
				d.id,
				d.nome,
				d.gerente_id,
				c.nome as gerente_nome,
				d.departamento_superior_id,
				dt.level + 1
			FROM departamentos d
			LEFT JOIN colaboradores c ON d.gerente_id = c.id
			INNER JOIN department_tree dt ON d.departamento_superior_id = dt.id
			WHERE d.deleted_at IS NULL
		)
		SELECT * FROM department_tree ORDER BY level, nome
	`

	var hierarchy []DeptHierarchy
	if err := r.db.Raw(query, id).Scan(&hierarchy).Error; err != nil {
		return nil, err
	}

	if len(hierarchy) == 0 {
		return nil, app.Errorf(app.ENOTFOUND, "department not found")
	}

	deptMap := make(map[uuid.UUID]*entities.Departamento)
	var rootDept *entities.Departamento

	for _, h := range hierarchy {
		dept := &entities.Departamento{
			ID:                     h.ID,
			Nome:                   h.Nome,
			GerenteID:              h.GerenteID,
			DepartamentoSuperiorID: h.DepartamentoSuperiorID,
			Gerente: &entities.Colaborador{
				ID:   h.GerenteID,
				Nome: h.GerenteNome,
			},
			Subdepartamentos: []entities.Departamento{},
		}
		deptMap[h.ID] = dept

		if h.Level == 0 {
			rootDept = dept
		}
	}

	for _, h := range hierarchy {
		if h.DepartamentoSuperiorID != nil {
			if parent, ok := deptMap[*h.DepartamentoSuperiorID]; ok {
				parent.Subdepartamentos = append(parent.Subdepartamentos, *deptMap[h.ID])
			}
		}
	}

	if r.cache != nil && rootDept != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("departamento:hierarquia:%s", id.String())
		if err := cache.SetJSON(ctx, r.cache, cacheKey, rootDept); err != nil {
			r.logger.Warn("Failed to cache department hierarchy", zap.Error(err), zap.String("key", cacheKey))
		}
	}

	return rootDept, nil
}

func (r *departamentoRepository) FindByFilters(filters map[string]interface{}, limit, offset int) ([]entities.Departamento, int64, error) {
	var departamentos []entities.Departamento
	var total int64

	query := r.db.Model(&entities.Departamento{})

	if nome, ok := filters["nome"].(string); ok && nome != "" {
		query = query.Where("nome ILIKE ?", "%"+nome+"%")
	}
	if gerenteNome, ok := filters["gerente_nome"].(string); ok && gerenteNome != "" {
		query = query.Joins("JOIN colaboradores ON colaboradores.id = departamentos.gerente_id").
			Where("colaboradores.nome ILIKE ?", "%"+gerenteNome+"%")
	}
	if parentID, ok := filters["departamento_superior_id"].(string); ok && parentID != "" {
		if id, err := uuid.Parse(parentID); err == nil {
			query = query.Where("departamento_superior_id = ?", id)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := query.Preload("Gerente").Preload("DepartamentoSuperior").
		Limit(limit).Offset(offset).Find(&departamentos)
	return departamentos, total, result.Error
}

func (r *departamentoRepository) FindSubDepartments(parentID uuid.UUID) ([]entities.Departamento, error) {
	var departamentos []entities.Departamento
	result := r.db.Where("departamento_superior_id = ?", parentID).Find(&departamentos)
	return departamentos, result.Error
}

func (r *departamentoRepository) FindAllSubDepartmentIDs(parentID uuid.UUID) ([]uuid.UUID, error) {
	if r.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("departamento:subdept_ids:%s", parentID.String())

		var ids []uuid.UUID
		err := cache.GetJSON(ctx, r.cache, cacheKey, &ids)
		if err == nil {
			return ids, nil
		}

		if err != redis.Nil {
			r.logger.Warn("Cache error", zap.Error(err), zap.String("key", cacheKey))
		}
	}

	query := `
		WITH RECURSIVE subdepartments AS (
			SELECT id FROM departamentos 
			WHERE id = ? AND deleted_at IS NULL
			
			UNION ALL
			
			SELECT d.id 
			FROM departamentos d
			INNER JOIN subdepartments sd ON d.departamento_superior_id = sd.id
			WHERE d.deleted_at IS NULL
		)
		SELECT id FROM subdepartments
	`

	var ids []uuid.UUID
	if err := r.db.Raw(query, parentID).Scan(&ids).Error; err != nil {
		return nil, err
	}

	if r.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("departamento:subdept_ids:%s", parentID.String())
		if err := cache.SetJSON(ctx, r.cache, cacheKey, ids); err != nil {
			r.logger.Warn("Failed to cache subdepartment IDs", zap.Error(err), zap.String("key", cacheKey))
		}
	}

	return ids, nil
}

func (r *departamentoRepository) HasCycle(departmentID, parentID uuid.UUID) (bool, error) {
	if departmentID == parentID {
		return true, nil
	}

	visited := make(map[uuid.UUID]bool)
	current := parentID

	for current != uuid.Nil {
		if current == departmentID {
			return true, nil
		}

		if visited[current] {
			return true, nil
		}
		visited[current] = true

		var dept entities.Departamento
		if err := r.db.First(&dept, "id = ?", current).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}

		if dept.DepartamentoSuperiorID == nil {
			break
		}
		current = *dept.DepartamentoSuperiorID
	}

	return false, nil
}

func (r *departamentoRepository) Create(departamento *entities.Departamento) error {
	err := r.db.Create(departamento).Error
	if err != nil {
		return err
	}

	if r.cache != nil {
		ctx := context.Background()
		patterns := []string{
			fmt.Sprintf("departamento:hierarquia:%s", departamento.ID.String()),
			fmt.Sprintf("departamento:subdept_ids:%s", departamento.ID.String()),
		}
		if departamento.DepartamentoSuperiorID != nil {
			patterns = append(patterns,
				fmt.Sprintf("departamento:hierarquia:%s", departamento.DepartamentoSuperiorID.String()),
				fmt.Sprintf("departamento:subdept_ids:%s", departamento.DepartamentoSuperiorID.String()),
			)
		}
		for _, pattern := range patterns {
			if err := r.cache.Del(ctx, pattern); err != nil {
				r.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("pattern", pattern))
			}
		}
	}

	return nil
}

func (r *departamentoRepository) Update(departamento *entities.Departamento) error {
	err := r.db.Model(departamento).Updates(departamento).Error
	if err != nil {
		return err
	}

	if r.cache != nil {
		ctx := context.Background()
		patterns := []string{
			fmt.Sprintf("departamento:hierarquia:%s", departamento.ID.String()),
			fmt.Sprintf("departamento:subdept_ids:%s", departamento.ID.String()),
		}
		if departamento.DepartamentoSuperiorID != nil {
			patterns = append(patterns,
				fmt.Sprintf("departamento:hierarquia:%s", departamento.DepartamentoSuperiorID.String()),
				fmt.Sprintf("departamento:subdept_ids:%s", departamento.DepartamentoSuperiorID.String()),
			)
		}
		for _, pattern := range patterns {
			if err := r.cache.Del(ctx, pattern); err != nil {
				r.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("pattern", pattern))
			}
		}
	}

	return nil
}

func (r *departamentoRepository) Delete(id uuid.UUID) error {
	dept, err := r.FindByID(id)
	if err != nil {
		return err
	}

	result := r.db.Delete(&entities.Departamento{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return app.Errorf(app.ENOTFOUND, "department not found")
	}

	if r.cache != nil {
		ctx := context.Background()
		patterns := []string{
			fmt.Sprintf("departamento:hierarquia:%s", id.String()),
			fmt.Sprintf("departamento:subdept_ids:%s", id.String()),
		}
		if dept.DepartamentoSuperiorID != nil {
			patterns = append(patterns,
				fmt.Sprintf("departamento:hierarquia:%s", dept.DepartamentoSuperiorID.String()),
				fmt.Sprintf("departamento:subdept_ids:%s", dept.DepartamentoSuperiorID.String()),
			)
		}
		for _, pattern := range patterns {
			if err := r.cache.Del(ctx, pattern); err != nil {
				r.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("pattern", pattern))
			}
		}
	}

	return nil
}

func (r *departamentoRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.Departamento{}).Count(&count)
	return count, result.Error
}

package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type cachedDepartamentoRepository struct {
	repo  DepartamentoRepository
	cache cache.Cache
}

func NewCachedDepartamentoRepository(repo DepartamentoRepository, cacheClient cache.Cache) DepartamentoRepository {
	return &cachedDepartamentoRepository{
		repo:  repo,
		cache: cacheClient,
	}
}

func (r *cachedDepartamentoRepository) FindAll(limit, offset int) ([]entities.Departamento, int64, error) {
	return r.repo.FindAll(limit, offset)
}

func (r *cachedDepartamentoRepository) FindByID(id uuid.UUID) (*entities.Departamento, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:%s", id.String())

	var departamento entities.Departamento
	err := cache.GetJSON(ctx, r.cache, cacheKey, &departamento)
	if err == nil {
		return &departamento, nil
	}

	if err != redis.Nil {
		fmt.Printf("Erro no cache: %v\n", err)
	}

	dept, err := r.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := cache.SetJSON(ctx, r.cache, cacheKey, dept); err != nil {
		fmt.Printf("Falha ao cachear departamento: %v\n", err)
	}

	return dept, nil
}

func (r *cachedDepartamentoRepository) FindByIDWithHierarchy(id uuid.UUID) (*entities.Departamento, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:hierarquia:%s", id.String())

	var departamento entities.Departamento
	err := cache.GetJSON(ctx, r.cache, cacheKey, &departamento)
	if err == nil {
		return &departamento, nil
	}

	if err != redis.Nil {
		fmt.Printf("Erro no cache: %v\n", err)
	}

	dept, err := r.repo.FindByIDWithHierarchy(id)
	if err != nil {
		return nil, err
	}

	if err := cache.SetJSON(ctx, r.cache, cacheKey, dept); err != nil {
		fmt.Printf("Falha ao cachear hierarquia do departamento: %v\n", err)
	}

	return dept, nil
}

func (r *cachedDepartamentoRepository) FindByFilters(filters map[string]interface{}, limit, offset int) ([]entities.Departamento, int64, error) {
	return r.repo.FindByFilters(filters, limit, offset)
}

func (r *cachedDepartamentoRepository) FindSubDepartments(parentID uuid.UUID) ([]entities.Departamento, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:subdepts:%s", parentID.String())

	var departamentos []entities.Departamento
	err := cache.GetJSON(ctx, r.cache, cacheKey, &departamentos)
	if err == nil {
		return departamentos, nil
	}

	if err != redis.Nil {
		fmt.Printf("Erro no cache: %v\n", err)
	}

	depts, err := r.repo.FindSubDepartments(parentID)
	if err != nil {
		return nil, err
	}

	if err := cache.SetJSON(ctx, r.cache, cacheKey, depts); err != nil {
		fmt.Printf("Falha ao cachear subdepartamentos: %v\n", err)
	}

	return depts, nil
}

func (r *cachedDepartamentoRepository) FindAllSubDepartmentIDs(parentID uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:subdept_ids:%s", parentID.String())

	var ids []uuid.UUID
	err := cache.GetJSON(ctx, r.cache, cacheKey, &ids)
	if err == nil {
		return ids, nil
	}

	if err != redis.Nil {
		fmt.Printf("Erro no cache: %v\n", err)
	}

	ids, err = r.repo.FindAllSubDepartmentIDs(parentID)
	if err != nil {
		return nil, err
	}

	if err := cache.SetJSON(ctx, r.cache, cacheKey, ids); err != nil {
		fmt.Printf("Falha ao cachear IDs dos subdepartamentos: %v\n", err)
	}

	return ids, nil
}

func (r *cachedDepartamentoRepository) HasCycle(departmentID, parentID uuid.UUID) (bool, error) {
	return r.repo.HasCycle(departmentID, parentID)
}

func (r *cachedDepartamentoRepository) Create(departamento *entities.Departamento) error {
	err := r.repo.Create(departamento)
	if err != nil {
		return err
	}

	r.invalidateDepartamentoCache(departamento.ID)
	if departamento.DepartamentoSuperiorID != nil {
		r.invalidateDepartamentoCache(*departamento.DepartamentoSuperiorID)
	}

	return nil
}

func (r *cachedDepartamentoRepository) Update(departamento *entities.Departamento) error {
	err := r.repo.Update(departamento)
	if err != nil {
		return err
	}

	r.invalidateDepartamentoCache(departamento.ID)
	if departamento.DepartamentoSuperiorID != nil {
		r.invalidateDepartamentoCache(*departamento.DepartamentoSuperiorID)
	}

	return nil
}

func (r *cachedDepartamentoRepository) Delete(id uuid.UUID) error {
	dept, err := r.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = r.repo.Delete(id)
	if err != nil {
		return err
	}

	r.invalidateDepartamentoCache(id)
	if dept.DepartamentoSuperiorID != nil {
		r.invalidateDepartamentoCache(*dept.DepartamentoSuperiorID)
	}

	return nil
}

func (r *cachedDepartamentoRepository) Count() (int64, error) {
	return r.repo.Count()
}

func (r *cachedDepartamentoRepository) invalidateDepartamentoCache(id uuid.UUID) {
	ctx := context.Background()
	patterns := []string{
		fmt.Sprintf("departamento:%s", id.String()),
		fmt.Sprintf("departamento:hierarquia:%s", id.String()),
		fmt.Sprintf("departamento:subdepts:%s", id.String()),
		fmt.Sprintf("departamento:subdept_ids:%s", id.String()),
	}

	for _, pattern := range patterns {
		if err := r.cache.Del(ctx, pattern); err != nil {
			fmt.Printf("Falha ao invalidar cache para %s: %v\n", pattern, err)
		}
	}

	if err := r.cache.DelPattern(ctx, fmt.Sprintf("*%s*", id.String())); err != nil {
		fmt.Printf("Falha ao invalidar padrão de cache: %v\n", err)
	}
}

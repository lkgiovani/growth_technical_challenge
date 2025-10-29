package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type cachedColaboradorRepository struct {
	repo  ColaboradorRepository
	cache cache.Cache
}

func NewCachedColaboradorRepository(repo ColaboradorRepository, cacheClient cache.Cache) ColaboradorRepository {
	return &cachedColaboradorRepository{
		repo:  repo,
		cache: cacheClient,
	}
}

func (r *cachedColaboradorRepository) FindAll(limit, offset int) ([]entities.Colaborador, int64, error) {
	return r.repo.FindAll(limit, offset)
}

func (r *cachedColaboradorRepository) FindByID(id uuid.UUID) (*entities.Colaborador, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("colaborador:%s", id.String())

	var colaborador entities.Colaborador
	err := cache.GetJSON(ctx, r.cache, cacheKey, &colaborador)
	if err == nil {
		return &colaborador, nil
	}

	if err != redis.Nil {
		fmt.Printf("Erro no cache: %v\n", err)
	}

	col, err := r.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := cache.SetJSON(ctx, r.cache, cacheKey, col); err != nil {
		fmt.Printf("Falha ao cachear colaborador: %v\n", err)
	}

	return col, nil
}

func (r *cachedColaboradorRepository) FindByCPF(cpf string) (*entities.Colaborador, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("colaborador:cpf:%s", cpf)

	var colaborador entities.Colaborador
	err := cache.GetJSON(ctx, r.cache, cacheKey, &colaborador)
	if err == nil {
		return &colaborador, nil
	}

	if err != redis.Nil {
		fmt.Printf("Erro no cache: %v\n", err)
	}

	col, err := r.repo.FindByCPF(cpf)
	if err != nil {
		return nil, err
	}

	if col != nil {
		if err := cache.SetJSON(ctx, r.cache, cacheKey, col); err != nil {
			fmt.Printf("Falha ao cachear colaborador por CPF: %v\n", err)
		}
	}

	return col, nil
}

func (r *cachedColaboradorRepository) FindByRG(rg string) (*entities.Colaborador, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("colaborador:rg:%s", rg)

	var colaborador entities.Colaborador
	err := cache.GetJSON(ctx, r.cache, cacheKey, &colaborador)
	if err == nil {
		return &colaborador, nil
	}

	if err != redis.Nil {
		fmt.Printf("Erro no cache: %v\n", err)
	}

	col, err := r.repo.FindByRG(rg)
	if err != nil {
		return nil, err
	}

	if col != nil {
		if err := cache.SetJSON(ctx, r.cache, cacheKey, col); err != nil {
			fmt.Printf("Falha ao cachear colaborador por RG: %v\n", err)
		}
	}

	return col, nil
}

func (r *cachedColaboradorRepository) FindByFilters(filters map[string]interface{}, limit, offset int) ([]entities.Colaborador, int64, error) {
	return r.repo.FindByFilters(filters, limit, offset)
}

func (r *cachedColaboradorRepository) FindByDepartmentIDs(departmentIDs []uuid.UUID) ([]entities.Colaborador, error) {
	ctx := context.Background()

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

		cols, err := r.repo.FindByDepartmentIDs(departmentIDs)
		if err != nil {
			return nil, err
		}

		if err := cache.SetJSON(ctx, r.cache, cacheKey, cols); err != nil {
			fmt.Printf("Falha ao cachear colaboradores por departamento: %v\n", err)
		}

		return cols, nil
	}

	return r.repo.FindByDepartmentIDs(departmentIDs)
}

func (r *cachedColaboradorRepository) Create(colaborador *entities.Colaborador) error {
	err := r.repo.Create(colaborador)
	if err != nil {
		return err
	}

	r.invalidateColaboradorCache(colaborador)
	return nil
}

func (r *cachedColaboradorRepository) Update(colaborador *entities.Colaborador) error {
	err := r.repo.Update(colaborador)
	if err != nil {
		return err
	}

	r.invalidateColaboradorCache(colaborador)
	return nil
}

func (r *cachedColaboradorRepository) Delete(id uuid.UUID) error {
	col, err := r.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = r.repo.Delete(id)
	if err != nil {
		return err
	}

	r.invalidateColaboradorCache(col)
	return nil
}

func (r *cachedColaboradorRepository) Count() (int64, error) {
	return r.repo.Count()
}

func (r *cachedColaboradorRepository) invalidateColaboradorCache(colaborador *entities.Colaborador) {
	ctx := context.Background()
	patterns := []string{
		fmt.Sprintf("colaborador:%s", colaborador.ID.String()),
		fmt.Sprintf("colaborador:cpf:%s", colaborador.CPF),
		fmt.Sprintf("colaborador:dept:%s", colaborador.DepartamentoID.String()),
	}

	if colaborador.RG != nil && *colaborador.RG != "" {
		patterns = append(patterns, fmt.Sprintf("colaborador:rg:%s", *colaborador.RG))
	}

	for _, pattern := range patterns {
		if err := r.cache.Del(ctx, pattern); err != nil {
			fmt.Printf("Falha ao invalidar cache para %s: %v\n", pattern, err)
		}
	}
}

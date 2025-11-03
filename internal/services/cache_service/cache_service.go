package cache_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lkgiovani/growth_technical_challenge/internal/domain/entities"
	"github.com/lkgiovani/growth_technical_challenge/pkg/cache"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type cacheService struct {
	cache  cache.Cache
	logger logger.Logger
}

func NewCacheService(cacheClient cache.Cache, log logger.Logger) CacheService {
	return &cacheService{
		cache:  cacheClient,
		logger: log,
	}
}

func (s *cacheService) GetColaboradoresByDepartment(departmentID uuid.UUID) ([]entities.Colaborador, error) {
	if s.cache == nil {
		return nil, redis.Nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("colaborador:dept:%s", departmentID.String())

	var colaboradores []entities.Colaborador
	err := cache.GetJSON(ctx, s.cache, cacheKey, &colaboradores)
	if err == nil {
		s.logger.Debug("Cache hit for colaboradores by department", zap.String("departmentID", departmentID.String()))
		return colaboradores, nil
	}

	if err != redis.Nil {
		s.logger.Warn("Cache error", zap.Error(err), zap.String("key", cacheKey))
	}

	return nil, err
}

func (s *cacheService) SetColaboradoresByDepartment(departmentID uuid.UUID, colaboradores []entities.Colaborador) error {
	if s.cache == nil {
		return nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("colaborador:dept:%s", departmentID.String())

	if err := cache.SetJSON(ctx, s.cache, cacheKey, colaboradores); err != nil {
		s.logger.Warn("Failed to cache colaboradores by department", zap.Error(err), zap.String("key", cacheKey))
		return err
	}

	return nil
}

func (s *cacheService) InvalidateColaboradoresByDepartment(departmentID uuid.UUID) error {
	if s.cache == nil {
		return nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("colaborador:dept:%s", departmentID.String())

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		s.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("key", cacheKey))
		return err
	}

	return nil
}

func (s *cacheService) GetDepartmentHierarchy(departmentID uuid.UUID) (*entities.Departamento, error) {
	if s.cache == nil {
		return nil, redis.Nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:hierarquia:%s", departmentID.String())

	var departamento entities.Departamento
	err := cache.GetJSON(ctx, s.cache, cacheKey, &departamento)
	if err == nil {
		s.logger.Debug("Cache hit for department hierarchy", zap.String("departmentID", departmentID.String()))
		return &departamento, nil
	}

	if err != redis.Nil {
		s.logger.Warn("Cache error", zap.Error(err), zap.String("key", cacheKey))
	}

	return nil, err
}

func (s *cacheService) SetDepartmentHierarchy(departmentID uuid.UUID, department *entities.Departamento) error {
	if s.cache == nil {
		return nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:hierarquia:%s", departmentID.String())

	if err := cache.SetJSON(ctx, s.cache, cacheKey, department); err != nil {
		s.logger.Warn("Failed to cache department hierarchy", zap.Error(err), zap.String("key", cacheKey))
		return err
	}

	return nil
}

func (s *cacheService) InvalidateDepartmentHierarchy(departmentID uuid.UUID) error {
	if s.cache == nil {
		return nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:hierarquia:%s", departmentID.String())

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		s.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("key", cacheKey))
		return err
	}

	return nil
}

func (s *cacheService) GetSubDepartmentIDs(parentID uuid.UUID) ([]uuid.UUID, error) {
	if s.cache == nil {
		return nil, redis.Nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:subdept_ids:%s", parentID.String())

	var ids []uuid.UUID
	err := cache.GetJSON(ctx, s.cache, cacheKey, &ids)
	if err == nil {
		s.logger.Debug("Cache hit for subdepartment IDs", zap.String("parentID", parentID.String()))
		return ids, nil
	}

	if err != redis.Nil {
		s.logger.Warn("Cache error", zap.Error(err), zap.String("key", cacheKey))
	}

	return nil, err
}

func (s *cacheService) SetSubDepartmentIDs(parentID uuid.UUID, ids []uuid.UUID) error {
	if s.cache == nil {
		return nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:subdept_ids:%s", parentID.String())

	if err := cache.SetJSON(ctx, s.cache, cacheKey, ids); err != nil {
		s.logger.Warn("Failed to cache subdepartment IDs", zap.Error(err), zap.String("key", cacheKey))
		return err
	}

	return nil
}

func (s *cacheService) InvalidateSubDepartmentIDs(parentID uuid.UUID) error {
	if s.cache == nil {
		return nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("departamento:subdept_ids:%s", parentID.String())

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		s.logger.Warn("Failed to invalidate cache", zap.Error(err), zap.String("key", cacheKey))
		return err
	}

	return nil
}

func (s *cacheService) InvalidateAllDepartmentCache(departmentID uuid.UUID, parentID *uuid.UUID) error {
	if s.cache == nil {
		return nil
	}

	if err := s.InvalidateDepartmentHierarchy(departmentID); err != nil {
		return err
	}

	if err := s.InvalidateSubDepartmentIDs(departmentID); err != nil {
		return err
	}

	if parentID != nil {
		if err := s.InvalidateDepartmentHierarchy(*parentID); err != nil {
			return err
		}

		if err := s.InvalidateSubDepartmentIDs(*parentID); err != nil {
			return err
		}
	}

	return nil
}

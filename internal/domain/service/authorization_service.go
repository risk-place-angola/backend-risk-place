package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

const (
	cacheExpiration = 5 * time.Minute
	cleanupInterval = 10 * time.Minute
)

type cacheEntry struct {
	permissions []model.Permission
	expiresAt   time.Time
}

type AuthorizationService struct {
	permRepo repository.PermissionRepository
	cache    sync.Map
}

func NewAuthorizationService(permRepo repository.PermissionRepository) *AuthorizationService {
	svc := &AuthorizationService{
		permRepo: permRepo,
	}
	go svc.cleanupExpiredCache()
	return svc
}

func (s *AuthorizationService) HasPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	permissionCode := resource + ":" + action

	permissions, err := s.getUserPermissionsWithCache(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm.Code == permissionCode {
			return true, nil
		}
	}
	return false, nil
}

func (s *AuthorizationService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error) {
	return s.getUserPermissionsWithCache(ctx, userID)
}

func (s *AuthorizationService) InvalidateUserCache(userID uuid.UUID) {
	s.cache.Delete(userID.String())
}

func (s *AuthorizationService) getUserPermissionsWithCache(ctx context.Context, userID uuid.UUID) ([]model.Permission, error) {
	key := userID.String()

	if val, ok := s.cache.Load(key); ok {
		entry, ok := val.(cacheEntry)
		switch {
		case !ok:
			s.cache.Delete(key)
		case time.Now().Before(entry.expiresAt):
			return entry.permissions, nil
		default:
			s.cache.Delete(key)
		}
	}

	permissions, err := s.permRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	s.cache.Store(key, cacheEntry{
		permissions: permissions,
		expiresAt:   time.Now().Add(cacheExpiration),
	})

	return permissions, nil
}

func (s *AuthorizationService) cleanupExpiredCache() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		s.cache.Range(func(key, value interface{}) bool {
			entry, ok := value.(cacheEntry)
			if !ok {
				s.cache.Delete(key)
				return true
			}
			if now.After(entry.expiresAt) {
				s.cache.Delete(key)
				slog.Debug("cleaned expired cache entry", "key", key)
			}
			return true
		})
	}
}

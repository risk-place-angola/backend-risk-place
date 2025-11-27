package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]model.Permission, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) HasPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error) {
	args := m.Called(ctx, userID, permissionCode)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]model.Permission, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func TestAuthorizationService_HasPermission(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	authzService := NewAuthorizationService(mockRepo)

	userID := uuid.New()
	ctx := context.Background()

	permissions := []model.Permission{
		{
			ID:       uuid.New(),
			Resource: "risk_type",
			Action:   "manage",
			Code:     "risk_type:manage",
		},
		{
			ID:       uuid.New(),
			Resource: "report",
			Action:   "read",
			Code:     "report:read",
		},
	}

	mockRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil)

	hasPermission, err := authzService.HasPermission(ctx, userID, "risk_type", "manage")
	assert.NoError(t, err)
	assert.True(t, hasPermission)

	hasPermission, err = authzService.HasPermission(ctx, userID, "risk_type", "delete")
	assert.NoError(t, err)
	assert.False(t, hasPermission)

	mockRepo.AssertExpectations(t)
}

func TestAuthorizationService_Cache(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	authzService := NewAuthorizationService(mockRepo)

	userID := uuid.New()
	ctx := context.Background()

	permissions := []model.Permission{
		{
			ID:       uuid.New(),
			Resource: "risk_type",
			Action:   "manage",
			Code:     "risk_type:manage",
		},
	}

	mockRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil).Once()

	_, err := authzService.HasPermission(ctx, userID, "risk_type", "manage")
	assert.NoError(t, err)

	_, err = authzService.HasPermission(ctx, userID, "risk_type", "manage")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthorizationService_InvalidateCache(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	authzService := NewAuthorizationService(mockRepo)

	userID := uuid.New()
	ctx := context.Background()

	permissions := []model.Permission{
		{
			ID:       uuid.New(),
			Resource: "risk_type",
			Action:   "manage",
			Code:     "risk_type:manage",
		},
	}

	mockRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil).Times(2)

	_, err := authzService.HasPermission(ctx, userID, "risk_type", "manage")
	assert.NoError(t, err)

	authzService.InvalidateUserCache(userID)

	_, err = authzService.HasPermission(ctx, userID, "risk_type", "manage")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthorizationService_CacheExpiration(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	authzService := NewAuthorizationService(mockRepo)

	userID := uuid.New()
	ctx := context.Background()

	permissions := []model.Permission{
		{
			ID:       uuid.New(),
			Resource: "risk_type",
			Action:   "manage",
			Code:     "risk_type:manage",
		},
	}

	mockRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil).Once()

	perms, err := authzService.GetUserPermissions(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, perms, 1)

	time.Sleep(100 * time.Millisecond)

	perms, err = authzService.GetUserPermissions(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, perms, 1)

	mockRepo.AssertExpectations(t)
}

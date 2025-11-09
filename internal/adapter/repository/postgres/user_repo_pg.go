package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type userRepoPG struct {
	q sqlc.Querier
}

// ListDeviceTokensByUserIDs implements repository.UserRepository.
func (u *userRepoPG) ListDeviceTokensByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]string, error) {
	rows, err := u.q.ListDeviceTokensByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	var tokens = make([]string, 0, len(rows))
	for _, row := range rows {
		if row.DeviceFcmToken.Valid {
			tokens = append(tokens, row.DeviceFcmToken.String)
		}
	}

	return tokens, nil
}

// ListAllDeviceTokensExceptUser implements repository.UserRepository.
func (u *userRepoPG) ListAllDeviceTokensExceptUser(ctx context.Context, excludeUserID uuid.UUID) ([]string, error) {
	rows, err := u.q.ListAllDeviceTokensExceptUser(ctx, excludeUserID)
	if err != nil {
		return nil, err
	}

	var tokens = make([]string, 0, len(rows))
	for _, row := range rows {
		if row.DeviceFcmToken.Valid {
			tokens = append(tokens, row.DeviceFcmToken.String)
		}
	}

	return tokens, nil
}

// UpdateUserDeviceInfo implements repository.UserRepository.
func (u *userRepoPG) UpdateUserDeviceInfo(ctx context.Context, userID uuid.UUID, fcmToken string, language string) error {
	return u.q.UpdateUserDeviceInfo(ctx, sqlc.UpdateUserDeviceInfoParams{
		ID:             userID,
		DeviceFcmToken: sql.NullString{String: fcmToken, Valid: fcmToken != ""},
		DeviceLanguage: sql.NullString{String: language, Valid: language != ""},
	})
}

func (u *userRepoPG) Save(ctx context.Context, entity *model.User) error {
	return u.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:       entity.ID,
		Name:     entity.Name,
		Email:    entity.Email,
		Password: entity.Password,
		Phone:    sql.NullString{String: entity.Phone, Valid: entity.Phone != ""},
	})
}

func (u *userRepoPG) Update(ctx context.Context, entity *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepoPG) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepoPG) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	userRow, err := u.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       userRow.ID,
		Name:     userRow.Name,
		Email:    userRow.Email,
		Password: userRow.Password,
		Phone:    userRow.Phone.String,
		Nif:      userRow.Nif.String,
		Address: model.Address{
			Country:      userRow.Country.String,
			Province:     userRow.Province.String,
			Municipality: userRow.Municipality.String,
			Neighborhood: userRow.Neighborhood.String,
			ZipCode:      userRow.ZipCode.String,
		},
		CreatedAt: userRow.CreatedAt.Time,
		UpdatedAt: userRow.UpdatedAt.Time,
	}

	return user, nil
}

func (u *userRepoPG) FindAll(ctx context.Context) ([]*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepoPG) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	userRow, err := u.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       userRow.ID,
		Name:     userRow.Name,
		Email:    userRow.Email,
		Password: userRow.Password,
		Phone:    userRow.Phone.String,
		Address: model.Address{
			Country:      userRow.Country.String,
			Province:     userRow.Province.String,
			Municipality: userRow.Municipality.String,
			Neighborhood: userRow.Neighborhood.String,
			ZipCode:      userRow.ZipCode.String,
		},
		CreatedAt: userRow.CreatedAt.Time,
		UpdatedAt: userRow.UpdatedAt.Time,
	}

	return user, nil
}

func (u *userRepoPG) AddCodeToUser(ctx context.Context, userID uuid.UUID, code string, expiration time.Time) error {
	return u.q.AddCodeToUser(ctx, sqlc.AddCodeToUserParams{
		ID:                         userID,
		EmailVerificationCode:      sql.NullString{String: code, Valid: code != ""},
		EmailVerificationExpiresAt: sql.NullTime{Time: expiration, Valid: !expiration.IsZero()},
	})
}

func (u *userRepoPG) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	return u.q.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:       userID,
		Password: newPassword,
	})
}

func (u *userRepoPG) UserHasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepoPG(db *sql.DB) repository.UserRepository {
	return &userRepoPG{
		q: sqlc.New(db),
	}
}

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
	q  sqlc.Querier
	db *sql.DB
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
	panic("implement me")
}

func (u *userRepoPG) Delete(ctx context.Context, id string) error {
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
	panic("implement me")
}

func (u *userRepoPG) UpdateSavedLocations(ctx context.Context, userID uuid.UUID, homeAddress, workAddress *model.SavedLocation) error {
	var homeNamePtr, homeAddressPtr *string
	var homeLatPtr, homeLonPtr *float64
	var workNamePtr, workAddressPtr *string
	var workLatPtr, workLonPtr *float64

	if homeAddress != nil {
		homeNamePtr = &homeAddress.Name
		homeAddressPtr = &homeAddress.Address
		homeLatPtr = &homeAddress.Latitude
		homeLonPtr = &homeAddress.Longitude
	}

	if workAddress != nil {
		workNamePtr = &workAddress.Name
		workAddressPtr = &workAddress.Address
		workLatPtr = &workAddress.Latitude
		workLonPtr = &workAddress.Longitude
	}

	return u.q.UpdateUserSavedLocations(ctx, sqlc.UpdateUserSavedLocationsParams{
		ID:                 userID,
		HomeAddressName:    sql.NullString{String: strOrEmpty(homeNamePtr), Valid: homeNamePtr != nil},
		HomeAddressAddress: sql.NullString{String: strOrEmpty(homeAddressPtr), Valid: homeAddressPtr != nil},
		HomeAddressLat:     sql.NullFloat64{Float64: floatOrZero(homeLatPtr), Valid: homeLatPtr != nil},
		HomeAddressLon:     sql.NullFloat64{Float64: floatOrZero(homeLonPtr), Valid: homeLonPtr != nil},
		WorkAddressName:    sql.NullString{String: strOrEmpty(workNamePtr), Valid: workNamePtr != nil},
		WorkAddressAddress: sql.NullString{String: strOrEmpty(workAddressPtr), Valid: workAddressPtr != nil},
		WorkAddressLat:     sql.NullFloat64{Float64: floatOrZero(workLatPtr), Valid: workLatPtr != nil},
		WorkAddressLon:     sql.NullFloat64{Float64: floatOrZero(workLonPtr), Valid: workLonPtr != nil},
	})
}

func strOrEmpty(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func floatOrZero(ptr *float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func (u *userRepoPG) UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, pushEnabled, smsEnabled bool) error {
	query := `
		UPDATE users 
		SET push_notification_enabled = $2, 
		    sms_notification_enabled = $3,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := u.db.ExecContext(ctx, query, userID, pushEnabled, smsEnabled)
	return err
}

func (u *userRepoPG) GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (pushEnabled, smsEnabled bool, err error) {
	query := `
		SELECT push_notification_enabled, sms_notification_enabled 
		FROM users 
		WHERE id = $1
	`
	err = u.db.QueryRowContext(ctx, query, userID).Scan(&pushEnabled, &smsEnabled)
	return
}

func (u *userRepoPG) GetUserLanguageAndPhone(ctx context.Context, userID uuid.UUID) (language, phone string, err error) {
	query := `
		SELECT COALESCE(device_language, 'pt'), COALESCE(phone, '') 
		FROM users 
		WHERE id = $1
	`
	err = u.db.QueryRowContext(ctx, query, userID).Scan(&language, &phone)
	return
}

func NewUserRepoPG(db *sql.DB) repository.UserRepository {
	return &userRepoPG{
		q:  sqlc.New(db),
		db: db,
	}
}

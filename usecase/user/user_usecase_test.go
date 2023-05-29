package user_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	user_usecase "github.com/risk-place-angola/backend-risk-place/usecase/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepository := mocks.NewMockUserRepository(ctrl)
	mockUserRepository.EXPECT().Save(gomock.Any()).Return(nil)

	userUseCase := user_usecase.NewUserUseCase(mockUserRepository)
	user, err := userUseCase.CreateUser(&user_usecase.CreateUserDTO{
		Name:     "Manuel Cunga",
		Phone:    "923456789",
		Email:    "manuel@gmail.com",
		Password: "909d04f13b3106f8633b38d",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Manuel Cunga", user.Name)

}

func TestFindAllUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.User{
		{
			ID:       "99bada49-09d0-4f13-b310-6f8633b38dfe",
			Name:     "Manuel Cunga",
			Email:    "manuel@gmail.com",
			Password: "909d04f13b3106f8633b38d",
		},
		{
			ID:       "99bada49-09d0-4f13-b310-6f8633b38dfe",
			Name:     "RRPL",
			Email:    "rrpl@gmail.com",
			Password: "909d04f13b3106f8633b38d",
		},
	}

	mockUserRepository := mocks.NewMockUserRepository(ctrl)

	mockUserRepository.EXPECT().FindAll().Return(data, nil)

	userUseCase := user_usecase.NewUserUseCase(mockUserRepository)
	user, err := userUseCase.FindAllUser()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(user))
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.User{
		ID:       "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		Name:     "omunga",
		Phone:    "923456789",
		Email:    "omunga@gmail.com",
		Password: "1234",
	}

	mockUserRepository := mocks.NewMockUserRepository(ctrl)
	mockUserRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockUserRepository.EXPECT().Update(gomock.Any()).Return(nil)

	updateuserDTO := &user_usecase.UpdateUserDTO{}
	updateuserDTO.Name = "Omunga plataforma"
	updateuserDTO.Phone = "923000000"
	updateuserDTO.Email = "omunga.io@gmail.com"
	updateuserDTO.Password = "12345"
	userUseCase := user_usecase.NewUserUseCase(mockUserRepository)
	user, err := userUseCase.UpdateUser("93247691-5c64-4c1f-a8ca-db5d76640ca9", updateuserDTO)
	assert.Nil(t, err)
	assert.Equal(t, "Omunga plataforma", user.Name)
}

func TestFindUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.User{
		ID:       "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		Name:     "Manuel",
		Email:    "manuel@gmail.com",
		Password: "asasasasasa",
	}

	mockUserRepository := mocks.NewMockUserRepository(ctrl)
	mockUserRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

	userUseCase := user_usecase.NewUserUseCase(mockUserRepository)
	user, err := userUseCase.FindUserByID("93247691-5c64-4c1f-a8ca-db5d76640ca9")
	assert.Nil(t, err)
	assert.Equal(t, "Manuel", user.Name)
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.User{
		ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
		Name:     "linkedin",
		Email:    "linkedin@gmail.com",
		Password: "12345",
	}

	mockUserpeRepository := mocks.NewMockUserRepository(ctrl)
	mockUserpeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockUserpeRepository.EXPECT().Delete(gomock.Any()).Return(nil)

	userUseCase := user_usecase.NewUserUseCase(mockUserpeRepository)
	err := userUseCase.RemoveUser("20dabe23-3541-455b-b64d-3191f2b2a303")
	assert.Nil(t, err)
}

func TestUserUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepository := mocks.NewMockUserRepository(ctrl)
	userUseCase := user_usecase.NewUserUseCase(mockUserRepository)

	t.Run("valid credentials should return token", func(t *testing.T) {
		// Mock FindByEmail to return a user with valid credentials
		email := "john.doe@example.com"
		password := "password"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockUserRepository.EXPECT().FindByEmail(email).Return(&entities.User{Email: email, Password: string(hashedPassword)}, nil)

		token, err := userUseCase.Login(&user_usecase.LoginDTO{
			Email:    email,
			Password: password,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("invalid credentials should return error", func(t *testing.T) {
		// Mock FindByEmail to return an error
		email := "jane.doe@example.com"
		password := "password"
		mockUserRepository.EXPECT().FindByEmail(email).Return(nil, fmt.Errorf("user not found"))

		token, err := userUseCase.Login(&user_usecase.LoginDTO{
			Email:    email,
			Password: password,
		})

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

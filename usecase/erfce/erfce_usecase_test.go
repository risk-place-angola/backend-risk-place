package erfce_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	erfce_usecase "github.com/risk-place-angola/backend-risk-place/usecase/erfce"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateErfce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
	mockErfceRepository.EXPECT().Save(gomock.Any()).Return(nil)

	erfceUseCase := erfce_usecase.NewErfceUseCase(mockErfceRepository)
	erfce, err := erfceUseCase.CreateErfce(&erfce_usecase.CreateErfceDTO{
		Name:     "any_name",
		Email:    "manuel@gmail.com",
		Password: "909d04f13b3106f8633b38d",
	})
	assert.Nil(t, err)
	assert.Equal(t, "any_name", erfce.Name)

}

func TestFindAllErfce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.Erfce{
		{
			ID:       "99bada49-09d0-4f13-b310-6f8633b38dfe",
			Name:     "any_name",
			Email:    "any_name@gmail.com",
			Password: "909d04f13b3106f8633b38d",
		},
		{
			ID:       "99bada49-09d0-4f13-b310-6f8633b38dfe",
			Name:     "RRPL",
			Email:    "rrpl@gmail.com",
			Password: "909d04f13b3106f8633b38d",
		},
	}

	mockErfceRepository := mocks.NewMockErfceRepository(ctrl)

	mockErfceRepository.EXPECT().FindAll().Return(data, nil)

	erfceUseCase := erfce_usecase.NewErfceUseCase(mockErfceRepository)
	erfce, err := erfceUseCase.FindAllErfce()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(erfce))
}

func TestUpdateErfce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.Erfce{
		ID:       "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		Name:     "omunga",
		Email:    "omunga@gmail.com",
		Password: "1234",
	}

	mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
	mockErfceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockErfceRepository.EXPECT().Update(gomock.Any()).Return(nil)

	updatErfceDTO := &erfce_usecase.UpdateErfceDTO{}
	updatErfceDTO.Name = "anonymous"
	updatErfceDTO.Email = "anonymous.io@gmail.com"
	updatErfceDTO.Password = "12345"
	erfceUseCase := erfce_usecase.NewErfceUseCase(mockErfceRepository)
	erfce, err := erfceUseCase.UpdateErfce("93247691-5c64-4c1f-a8ca-db5d76640ca9", updatErfceDTO)
	assert.Nil(t, err)
	assert.Equal(t, "anonymous", erfce.Name)
}

func TestFindErfceID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.Erfce{
		ID:       "93247691-5c64-4c1f-a8ca-db5d76640ca9",
		Name:     "Manuel",
		Email:    "manuel@gmail.com",
		Password: "asasasasasa",
	}

	mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
	mockErfceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

	erfceUseCase := erfce_usecase.NewErfceUseCase(mockErfceRepository)
	erfce, err := erfceUseCase.FindErfceByID("93247691-5c64-4c1f-a8ca-db5d76640ca9")
	assert.Nil(t, err)
	assert.Equal(t, "Manuel", erfce.Name)
}

func TestDeleteErfce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := &entities.Erfce{
		ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
		Name:     "youtube",
		Email:    "youtube@gmail.com",
		Password: "12345",
	}

	mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
	mockErfceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
	mockErfceRepository.EXPECT().Delete(gomock.Any()).Return(nil)

	erfceUseCase := erfce_usecase.NewErfceUseCase(mockErfceRepository)
	err := erfceUseCase.RemoveErfce("20dabe23-3541-455b-b64d-3191f2b2a303")
	assert.Nil(t, err)
}

func TestUErfceUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
	erfceUseCase := erfce_usecase.NewErfceUseCase(mockErfceRepository)

	t.Run("valid credentials should return token", func(t *testing.T) {
		// Mock FindByEmail to return a user with valid credentials
		email := "john.doe@example.com"
		password := "password"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockErfceRepository.EXPECT().FindByEmail(email).Return(&entities.Erfce{Email: email, Password: string(hashedPassword)}, nil)

		token, err := erfceUseCase.Login(&erfce_usecase.LoginDTO{
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
		mockErfceRepository.EXPECT().FindByEmail(email).Return(nil, fmt.Errorf("user not found"))

		token, err := erfceUseCase.Login(&erfce_usecase.LoginDTO{
			Email:    email,
			Password: password,
		})

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

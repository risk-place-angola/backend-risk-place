package user_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	user_usecase "github.com/risk-place-angola/backend-risk-place/usecase/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepository := mocks.NewMockUserRepository(ctrl)
	mockUserRepository.EXPECT().Save(gomock.Any()).Return(nil)

	userUseCase := user_usecase.NewUserUseCase(mockUserRepository)
	user, err := userUseCase.CreateUser(&user_usecase.CreateUserDTO{
		ID:       "99bada49-09d0-4f13-b310-6f8633b38dfe",
		Name:     "Manuel Cunga",
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

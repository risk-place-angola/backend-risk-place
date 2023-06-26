package user_controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	user_controller "github.com/risk-place-angola/backend-risk-place/infra/rest/user/controllers"
	"golang.org/x/crypto/bcrypt"

	account "github.com/risk-place-angola/backend-risk-place/usecase/user"
	"github.com/stretchr/testify/assert"
)

func TestUserControllers(t *testing.T) {
	t.Run("should return 201 when create a new user", func(t *testing.T) {

		e := echo.New()

		data := entities.User{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "any_name",
			Phone:    "923456789",
			Email:    "joe@gmail.com",
			Password: "a80e68232ecd",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("POST", "/infra/v1/user/create", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepository := mocks.NewMockUserRepository(ctrl)
		mockUserRepository.EXPECT().Save(gomock.Any()).Return(nil)

		userUseCase := account.NewUserUseCase(mockUserRepository)
		userController := user_controller.NewUserController(userUseCase)

		if assert.NoError(t, userController.UserCreateController(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code, "error status code != 201")
		}
	})

	t.Run("should return 200 when update a user ", func(t *testing.T) {
		e := echo.New()

		data := &entities.User{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "Github",
			Phone:    "923456789",
			Email:    "github@gmail.com",
			Password: "12345",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("PUT", "/infra/v1/user/:id", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/infra/v1/user/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserpeRepository := mocks.NewMockUserRepository(ctrl)
		mockUserpeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
		mockUserpeRepository.EXPECT().Update(gomock.Any()).Return(nil)

		userUseCase := account.NewUserUseCase(mockUserpeRepository)
		userController := user_controller.NewUserController(userUseCase)

		if assert.NoError(t, userController.UserUpdateController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

	t.Run("should return 200 when delete a user", func(t *testing.T) {
		e := echo.New()

		data := &entities.User{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "Github",
			Phone:    "923456789",
			Email:    "github@gmail.com",
			Password: "12345",
		}

		res := httptest.NewRequest("DELETE", "/infra/v1/user/:id", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/infra/v1/user/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserpeRepository := mocks.NewMockUserRepository(ctrl)
		mockUserpeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
		mockUserpeRepository.EXPECT().Delete(gomock.Any()).Return(nil)

		userUseCase := account.NewUserUseCase(mockUserpeRepository)
		userController := user_controller.NewUserController(userUseCase)

		if assert.NoError(t, userController.UserDeleteController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

	t.Run("should return 200 when find all user", func(t *testing.T) {
		e := echo.New()

		data := []*entities.User{
			{
				ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
				Name:     "Github",
				Phone:    "923456789",
				Email:    "github@gmail.com",
				Password: "12345",
			},
			{
				ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
				Name:     "linkedin",
				Email:    "linkedin@gmail.com",
				Password: "12345",
			},
		}

		res := httptest.NewRequest("GET", "/infra/v1/users", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetPath("/infra/v1/users")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserpeRepository := mocks.NewMockUserRepository(ctrl)
		mockUserpeRepository.EXPECT().FindAll().Return(data, nil)

		userUseCase := account.NewUserUseCase(mockUserpeRepository)
		userController := user_controller.NewUserController(userUseCase)

		if assert.NoError(t, userController.UserFindAllController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

	t.Run("should return 200 when find a user by id", func(t *testing.T) {
		e := echo.New()

		data := &entities.User{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "linkedin",
			Phone:    "923456789",
			Email:    "linkedin@gmail.com",
			Password: "12345",
		}

		res := httptest.NewRequest("GET", "/infra/v1/user/:id", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/infra/v1/user/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserpeRepository := mocks.NewMockUserRepository(ctrl)
		mockUserpeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

		userUseCase := account.NewUserUseCase(mockUserpeRepository)
		userController := user_controller.NewUserController(userUseCase)

		if assert.NoError(t, userController.UserFindByIdController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

	t.Run("should return 200 when login is successful", func(t *testing.T) {

		e := echo.New()

		email := "john.doe@example.com"
		password := "password"

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		dto := account.LoginDTO{
			Email:    email,
			Password: password,
		}

		jsonData, _ := json.Marshal(dto)

		res := httptest.NewRequest("POST", "/infra/v1/user/login", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepository := mocks.NewMockUserRepository(ctrl)
		mockUserRepository.EXPECT().FindByEmail(email).Return(&entities.User{Email: email, Password: string(hashedPassword)}, nil)

		userUseCase := account.NewUserUseCase(mockUserRepository)
		userController := user_controller.NewUserController(userUseCase)

		if assert.NoError(t, userController.UserLoginController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

}

package user_controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	user_controller "github.com/risk-place-angola/backend-risk-place/app/rest/user/controllers"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	account "github.com/risk-place-angola/backend-risk-place/usecase/user"
	"github.com/stretchr/testify/assert"
)

func TestUserControllers(t *testing.T) {
	t.Run("should return 201 when create a new user", func(t *testing.T) {

		e := echo.New()

		data := entities.User{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "any_name",
			Email:    "joe@gmail.com",
			Password: "a80e68232ecd",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("POST", "/api/v1/user/create", bytes.NewBuffer(jsonData))
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
			Email:    "github@gmail.com",
			Password: "12345",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("PUT", "/api/v1/user/:id", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/api/v1/user/:id")

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
			Email:    "github@gmail.com",
			Password: "12345",
		}

		res := httptest.NewRequest("DELETE", "/api/v1/user/:id", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/api/v1/user/:id")

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

		res := httptest.NewRequest("GET", "/api/v1/users", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetPath("/api/v1/users")

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
			Email:    "linkedin@gmail.com",
			Password: "12345",
		}

		res := httptest.NewRequest("GET", "/api/v1/user/:id", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/api/v1/user/:id")

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

	t.Run("should return 200 when login with valid credentials", func(t *testing.T) {
		e := echo.New()

		credentials := account.LoginDTO{
			Email:    "joe@gmail.com",
			Password: "valid_password",
		}

		jsonData, _ := json.Marshal(credentials)

		res := httptest.NewRequest("POST", "/api/v1/user/auth", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepository := mocks.NewMockUserRepository(ctrl)
		mockUserRepository.EXPECT().FindByEmail(credentials.Email).Return(&entities.User{
			ID:       "valid_id",
			Name:     "valid_name",
			Email:    "joe@gmail.com",
			Password: "valid_encrypted_password",
		}, nil)

		userUseCase := account.NewUserUseCase(mockUserRepository)
		userController := user_controller.NewUserController(userUseCase)

		if assert.NoError(t, userController.UserLoginController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}
	})

	t.Run("should return 401 when login with invalid credentials", func(t *testing.T) {
		e := echo.New()

		credentials := account.LoginDTO{
			Email:    "joe@gmail.com",
			Password: "invalid_password",
		}

		jsonData, _ := json.Marshal(credentials)

		res := httptest.NewRequest("POST", "/api/v1/user/auth", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepository := mocks.NewMockUserRepository(ctrl)

		mockUserRepository.EXPECT().FindByEmail(credentials.Email).Return(&entities.User{
			Email: "joe@gmail.com",
		}, nil)

		userUseCase := account.NewUserUseCase(mockUserRepository)
		userController := user_controller.NewUserController(userUseCase)

		if assert.NoError(t, userController.UserLoginController(ctx)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code, "error status code != 401")
		}
	})

}

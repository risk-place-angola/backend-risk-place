package erfce_controller_test

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
	erfce_controller "github.com/risk-place-angola/backend-risk-place/infra/rest/erfce/controllers"
	"golang.org/x/crypto/bcrypt"

	account "github.com/risk-place-angola/backend-risk-place/usecase/erfce"
	"github.com/stretchr/testify/assert"
)

func TestErfceControllers(t *testing.T) {
	t.Run("should return 201 when create a new Erfce user", func(t *testing.T) {

		e := echo.New()

		data := entities.Erfce{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "any_name",
			Email:    "joe@gmail.com",
			Password: "a80e68232ecd",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("POST", "/infra/v1/erfce/create", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
		mockErfceRepository.EXPECT().Save(gomock.Any()).Return(nil)

		erfceUseCase := account.NewErfceUseCase(mockErfceRepository)
		erfceController := erfce_controller.NewErfceController(erfceUseCase)

		if assert.NoError(t, erfceController.ErfceCreateController(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code, "error status code != 201")
		}
	})

	t.Run("should return 200 when update a erfce ", func(t *testing.T) {
		e := echo.New()

		data := &entities.Erfce{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "Github",
			Email:    "github@gmail.com",
			Password: "12345",
		}

		jsonData, _ := json.Marshal(data)

		res := httptest.NewRequest("PUT", "/infra/v1/erfce/:id", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/infra/v1/erfce/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
		mockErfceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
		mockErfceRepository.EXPECT().Update(gomock.Any()).Return(nil)

		erfceUseCase := account.NewErfceUseCase(mockErfceRepository)
		erfceController := erfce_controller.NewErfceController(erfceUseCase)

		if assert.NoError(t, erfceController.ErfceUpdateController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

	t.Run("should return 200 when delete a erfce user", func(t *testing.T) {
		e := echo.New()

		data := &entities.Erfce{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "Github",
			Email:    "github@gmail.com",
			Password: "12345",
		}

		res := httptest.NewRequest("DELETE", "/infra/v1/erfce/:id", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/infra/v1/erfce/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
		mockErfceRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)
		mockErfceRepository.EXPECT().Delete(gomock.Any()).Return(nil)

		erfceUseCase := account.NewErfceUseCase(mockErfceRepository)
		erfceController := erfce_controller.NewErfceController(erfceUseCase)

		if assert.NoError(t, erfceController.ErfceDeleteController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

	t.Run("should return 200 when find all erfce users", func(t *testing.T) {
		e := echo.New()

		data := []*entities.Erfce{
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

		res := httptest.NewRequest("GET", "/infra/v1/erfces", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetPath("/infra/v1/erfces")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
		mockErfceRepository.EXPECT().FindAll().Return(data, nil)

		erfceUseCase := account.NewErfceUseCase(mockErfceRepository)
		erfceController := erfce_controller.NewErfceController(erfceUseCase)

		if assert.NoError(t, erfceController.ErfceFindAllController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

	t.Run("should return 200 when find a erfce user by id", func(t *testing.T) {
		e := echo.New()

		data := &entities.Erfce{
			ID:       "0c1baa42-3909-4bdb-837f-a80e68232ecd",
			Name:     "linkedin",
			Email:    "linkedin@gmail.com",
			Password: "12345",
		}

		res := httptest.NewRequest("GET", "/infra/v1/erfce/:id", nil)
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("0c1baa42-3909-4bdb-837f-a80e68232ecd")
		ctx.SetPath("/infra/v1/erfce/:id")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErfcepeRepository := mocks.NewMockErfceRepository(ctrl)
		mockErfcepeRepository.EXPECT().FindByID(gomock.Any()).Return(data, nil)

		erfceUseCase := account.NewErfceUseCase(mockErfcepeRepository)
		erfceController := erfce_controller.NewErfceController(erfceUseCase)

		if assert.NoError(t, erfceController.ErfceFindByIdController(ctx)) {
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

		res := httptest.NewRequest("POST", "/infra/v1/erfce/login", bytes.NewBuffer(jsonData))
		res.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(res, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErfceRepository := mocks.NewMockErfceRepository(ctrl)
		mockErfceRepository.EXPECT().FindByEmail(email).Return(&entities.Erfce{Email: email, Password: string(hashedPassword)}, nil)

		erfceUseCase := account.NewErfceUseCase(mockErfceRepository)
		erfceController := erfce_controller.NewErfceController(erfceUseCase)

		if assert.NoError(t, erfceController.ErfceLoginController(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code, "error status code != 200")
		}

	})

}

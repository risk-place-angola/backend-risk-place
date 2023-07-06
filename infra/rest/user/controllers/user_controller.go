package user_controller

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/infra/rest"
	user_presenter "github.com/risk-place-angola/backend-risk-place/infra/rest/user/presenter"
	account "github.com/risk-place-angola/backend-risk-place/usecase/user"
)

type UserController interface {
	UserCreateController(ctx user_presenter.UserPresenterCTX) error
	UserUpdateController(ctx user_presenter.UserPresenterCTX) error
	UserFindAllController(ctx user_presenter.UserPresenterCTX) error
	UserFindByIdController(ctx user_presenter.UserPresenterCTX) error
	UserDeleteController(ctx user_presenter.UserPresenterCTX) error
	UserLoginController(ctx user_presenter.UserPresenterCTX) error
	FindAllUserWarningsController(ctx user_presenter.UserPresenterCTX) error
	FindWarningByUserIDController(ctx user_presenter.UserPresenterCTX) error
}

type UserControllerImpl struct {
	userUseCase account.UserUseCase
}

func NewUserController(userRepo account.UserUseCase) UserController {
	return &UserControllerImpl{
		userUseCase: userRepo,
	}
}

// @Summary Create User
// @Description Create User
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body account.CreateUserDTO true "User"
// @Success 201 {object} account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/user [post]
func (controller *UserControllerImpl) UserCreateController(ctx user_presenter.UserPresenterCTX) error {
	var user account.CreateUserDTO
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	userCreate, err := controller.userUseCase.CreateUser(&user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, userCreate)
}

// @Summary Find All User
// @Description Find All User
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} []account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/user [get]
func (controller *UserControllerImpl) UserFindAllController(ctx user_presenter.UserPresenterCTX) error {
	users, err := controller.userUseCase.FindAllUser()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, users)
}

// @Summary Find User By ID
// @Description Find User By ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/user/{id} [get]
func (controller *UserControllerImpl) UserFindByIdController(ctx user_presenter.UserPresenterCTX) error {
	id := ctx.Param("id")

	userId, err := controller.userUseCase.FindUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, userId)
}

// @Summary Update User
// @Description Update User
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body account.UpdateUserDTO true "User"
// @Success 200 {object} account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/user/{id} [put]
func (controller *UserControllerImpl) UserUpdateController(ctx user_presenter.UserPresenterCTX) error {
	id := ctx.Param("id")

	var user account.UpdateUserDTO
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	userUpdate, err := controller.userUseCase.UpdateUser(id, &user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, userUpdate)
}

// @Summary Delete User
// @Description Delete User
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} rest.SuccessResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/user/{id} [delete]
func (controller *UserControllerImpl) UserDeleteController(ctx user_presenter.UserPresenterCTX) error {
	id := ctx.Param("id")

	err := controller.userUseCase.RemoveUser(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, rest.SuccessResponse{Message: "User deleted successfully"})
}

// @Summary Login User
// @Description Login User
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body account.LoginDTO true "User"
// @Success 200 {object} account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/user/login [post]
func (controller *UserControllerImpl) UserLoginController(ctx user_presenter.UserPresenterCTX) error {
	var credentials account.LoginDTO
	if err := ctx.Bind(&credentials); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	data, err := controller.userUseCase.Login(&credentials)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}

// @Summary Find All User Warnings
// @Description Find All User Warnings
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} []account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/user/warning [get]
func (controller *UserControllerImpl) FindAllUserWarningsController(ctx user_presenter.UserPresenterCTX) error {
	warnings, err := controller.userUseCase.FindAllUserWarnings()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, warnings)
}

// @Summary Find User Warnings By ID
// @Description Find User Warnings By ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} []account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/user/warning/{id} [get]
func (controller *UserControllerImpl) FindWarningByUserIDController(ctx user_presenter.UserPresenterCTX) error {
	id := ctx.Param("id")

	warnings, err := controller.userUseCase.FindWarningByUserID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, warnings)
}

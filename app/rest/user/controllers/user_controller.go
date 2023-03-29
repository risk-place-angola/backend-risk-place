package user_controller

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/app/rest"
	user_presenter "github.com/risk-place-angola/backend-risk-place/app/rest/user/presenter"
	account "github.com/risk-place-angola/backend-risk-place/usecase/user"
)

type UserController interface {
	UserCreateController(ctx user_presenter.UserPresenterCTX) error
	UserUpdateController(ctx user_presenter.UserPresenterCTX) error
	UserFindAllController(ctx user_presenter.UserPresenterCTX) error
	UserFindByIdController(ctx user_presenter.UserPresenterCTX) error
	UserDeleteController(ctx user_presenter.UserPresenterCTX) error
	UserLoginController(ctx user_presenter.UserPresenterCTX) error
}

type UserControllerImpl struct {
	userUseCase account.UserUseCase
}

func NewUserController(userRepo account.UserUseCase) UserController {
	return &UserControllerImpl{
		userUseCase: userRepo,
	}
}

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

func (controller *UserControllerImpl) UserFindAllController(ctx user_presenter.UserPresenterCTX) error {
	users, err := controller.userUseCase.FindAllUser()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, users)
}

func (controller *UserControllerImpl) UserFindByIdController(ctx user_presenter.UserPresenterCTX) error {
	id := ctx.Param("id")

	userId, err := controller.userUseCase.FindUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, userId)
}

func (controller *UserControllerImpl) UserUpdateController(ctx user_presenter.UserPresenterCTX) error {
	id := ctx.Param("id")

	var user account.UpadateUserDTO
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	userUpdate, err := controller.userUseCase.UpdateUser(id, &user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, userUpdate)
}

func (controller *UserControllerImpl) UserDeleteController(ctx user_presenter.UserPresenterCTX) error {
	id := ctx.Param("id")

	err := controller.userUseCase.RemoveUser(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, rest.SuccessResponse{Message: "User deleted successfully"})
}

func (controller *UserControllerImpl) UserLoginController(ctx user_presenter.UserPresenterCTX) error {
	var credentials account.LoginDTO
	if err := ctx.Bind(&credentials); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	data, err := controller.userUseCase.UserLogin(&credentials)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)

}

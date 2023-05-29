package authjwt

import "github.com/labstack/echo/v4"

type IAuthService interface {
	Auth(ctx echo.Context) error
}

type AuthService struct {
	IAuthAPI IAuthAPI
}

type Data struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthService(auth IAuthAPI) IAuthService {
	return &AuthService{IAuthAPI: auth}
}

// Auth is Authentication service
// @Summary auth
// @Description
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param place body Data true "Auth"
// @Success 200 {object} Token
// @Failure 401 {object} string	"UnAuthorized"
// @Router /auth [post]
func (a *AuthService) Auth(ctx echo.Context) error {
	var data Data
	if err := ctx.Bind(&data); err != nil {
		return ctx.JSON(400, err.Error())
	}
	token, err := a.IAuthAPI.Auth(data.Username, data.Password)
	if err != nil {
		return ctx.JSON(400, err.Error())
	}

	return ctx.JSON(200, token)
}

// Auths is Authentication service
// @Summary auths
// @Description return all auths
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} []entities.Auth
// @Failure 401 {object} string
// @Router /auths [get]
func (a *AuthService) Auths(ctx echo.Context) error {
	auths, err := a.IAuthAPI.Auths()
	if err != nil {
		return ctx.JSON(400, err.Error())
	}
	return ctx.JSON(200, auths)
}

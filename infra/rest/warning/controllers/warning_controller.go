package warning_controllers

import (
	"encoding/json"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/warning/presenter"
	warning_usecase "github.com/risk-place-angola/backend-risk-place/usecase/warning"
	"github.com/risk-place-angola/backend-risk-place/util"
	"log"
)

type IWarningController interface {
	CreateWarning(ctx presenter.WarningPresenterCTX) error
	UpdateWarning(ctx presenter.WarningPresenterCTX) error
	FindAllWarning(ctx presenter.WarningPresenterCTX) error
	FindWarningByID(ctx presenter.WarningPresenterCTX) error
	RemoveWarning(ctx presenter.WarningPresenterCTX) error
}

type WarningControllerImpl struct {
	IWarningUseCase        warning_usecase.IWarningUseCase
	WebsocketClientManager *util.WebsocketClientManager
	Env                    *util.Env
}

func NewWarningController(warningUseCase warning_usecase.IWarningUseCase) IWarningController {
	return &WarningControllerImpl{
		IWarningUseCase: warningUseCase,
		Env:             util.LoadEnv(".env"),
	}
}

// CreateWarning godoc
// @Summary Create a warning
// @Description Create a warning
// @Tags Warning
// @Accept  json
// @Produce  json
// @Param reported_by formData string true "Reported by"
// @Param latitude formData string true "Latitude"
// @Param longitude formData string true "Longitude"
// @Param fact formData file true "Fact"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /api/v1/warning [post]
func (w WarningControllerImpl) CreateWarning(ctx presenter.WarningPresenterCTX) error {
	warningDTO := &warning_usecase.CreateWarningDTO{}
	warningDTO.ReportedBy = ctx.FormValue("reported_by")
	warningDTO.Latitude = ctx.FormValue("latitude")
	warningDTO.Longitude = ctx.FormValue("longitude")
	file, err := ctx.FormFile("fact")
	if err != nil {
		return ctx.JSON(400, err.Error())
	}
	log.Println("env", w.Env.REMOTEHOST)
	conn, errWS := util.WebsocketClientDialer(w.Env.REMOTEHOST, ctx)
	if errWS != nil {
		return ctx.JSON(errWS.Code, errWS.Message)
	}

	manage := util.NewWebsocketClientManager()
	client := &util.Websocket{
		ID:                     warningDTO.ReportedBy,
		Conn:                   conn,
		Send:                   make(chan []byte),
		WebsocketClientManager: manage,
	}

	go manage.Start()

	manage.Register <- client

	s3, err := util.UploadFileToS3(file)
	if err != nil {
		return ctx.JSON(400, err.Error())
	}
	warningDTO.Fact = s3.Location

	warningOutputDTO, err := w.IWarningUseCase.CreateWarning(warningDTO)
	if err != nil {
		return ctx.JSON(400, err)
	}

	warningOutputDTOBytes, err := json.Marshal(warningOutputDTO)
	if err != nil {
		return ctx.JSON(400, err)
	}
	go client.WebsocketClientWriteMessage(warningOutputDTOBytes)

	return nil
}

// UpdateWarning godoc
// @Summary Update a warning
// @Description Update a warning
// @Tags Warning
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param update_warning body warning_usecase.UpdateWarningDTO true "Update Warning"
// @Success 200 {object} warning_usecase.UpdateWarningDTO
// @Failure 400 {object} string
// @Router /api/v1/warning/{id} [put]
func (w WarningControllerImpl) UpdateWarning(ctx presenter.WarningPresenterCTX) error {
	var warningDTO warning_usecase.UpdateWarningDTO
	id := ctx.Param("id")
	if err := ctx.Bind(&warningDTO); err != nil {
		return ctx.JSON(400, err.Error())
	}

	conn, errWS := util.WebsocketClientDialer(w.Env.REMOTEHOST, ctx)
	if errWS != nil {
		return ctx.JSON(errWS.Code, errWS.Message)
	}

	manage := util.NewWebsocketClientManager()
	client := &util.Websocket{
		ID:                     id,
		Conn:                   conn,
		Send:                   make(chan []byte),
		WebsocketClientManager: manage,
	}

	go manage.Start()

	manage.Register <- client

	warningOutputDTO, err := w.IWarningUseCase.UpdateWarning(id, &warningDTO)
	if err != nil {
		return ctx.JSON(400, err.Error())
	}

	warningOutputDTOBytes, err := json.Marshal(warningOutputDTO)
	if err != nil {
		return ctx.JSON(400, err)
	}
	go client.WebsocketClientWriteMessage(warningOutputDTOBytes)

	return nil
}

// FindAllWarning godoc
// @Summary Find all warnings
// @Description Find all warnings
// @Tags Warning
// @Accept  json
// @Produce  json
// @Success 200 {object} []warning_usecase.DTO
// @Failure 400 {object} string
// @Router /api/v1/warning [get]
func (w WarningControllerImpl) FindAllWarning(ctx presenter.WarningPresenterCTX) error {
	warnings, err := w.IWarningUseCase.FindAllWarning()
	if err != nil {
		return ctx.JSON(400, err.Error())
	}
	return ctx.JSON(200, warnings)
}

// FindWarningByID godoc
// @Summary Find warning by ID
// @Description Find warning by ID
// @Tags Warning
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} warning_usecase.DTO
// @Failure 400 {object} string
// @Router /api/v1/warning/{id} [get]
func (w WarningControllerImpl) FindWarningByID(ctx presenter.WarningPresenterCTX) error {
	id := ctx.Param("id")
	warning, err := w.IWarningUseCase.FindWarningByID(id)
	if err != nil {
		return ctx.JSON(400, err.Error())
	}
	return ctx.JSON(200, warning)
}

// RemoveWarning godoc
// @Summary Remove warning
// @Description Remove warning
// @Tags Warning
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /api/v1/warning/{id} [delete]
func (w WarningControllerImpl) RemoveWarning(ctx presenter.WarningPresenterCTX) error {
	id := ctx.Param("id")
	err := w.IWarningUseCase.RemoveWarning(id)
	if err != nil {
		return ctx.JSON(400, err.Error())
	}
	return ctx.JSON(200, "Warning removed")
}

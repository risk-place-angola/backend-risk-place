package place_controller

import (
	"net/http"

	place_presenter "github.com/risk-place-angola/backend-risk-place/api/rest/place/presenter"
	place_usecase "github.com/risk-place-angola/backend-risk-place/usecase/place"
)

type PlaceController interface {
	PlaceCreateController(ctx place_presenter.PlacePresenterCTX) error
	PlaceFindByIdController(ctx place_presenter.PlacePresenterCTX) error
}

type PlaceControllerImpl struct {
	placeUseCase place_usecase.PlaceUseCase
}

func NewPlaceController(placeUseCase place_usecase.PlaceUseCase) PlaceController {
	return &PlaceControllerImpl{
		placeUseCase: placeUseCase,
	}
}

// @Summary Create Place
// @Description Create Place
// @Tags Place
// @Accept  json
// @Produce  json
// @Param place body place_usecase.CreatePlaceDTO true "Place"
// @Success 201 {object} place_usecase.PlaceDTO
// @Failure 500 {object} place_presenter.ErrorResponse
// @Router /api/v1/place [post]
func (controller *PlaceControllerImpl) PlaceCreateController(ctx place_presenter.PlacePresenterCTX) error {
	var place place_usecase.CreatePlaceDTO
	if err := ctx.Bind(&place); err != nil {
		return ctx.JSON(http.StatusBadRequest, place_presenter.ErrorResponse{Message: err.Error()})
	}

	placeCreate, err := controller.placeUseCase.CreatePlace(place)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, place_presenter.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, placeCreate)
}

// @Summary Find Place by ID
// @Description Find Place by ID
// @Tags Place
// @Accept  json
// @Produce  json
// @Param id path string true "Place ID"
// @Success 200 {object} place_usecase.PlaceDTO
// @Failure 500 {object} place_presenter.ErrorResponse
// @Router /api/v1/place/{id} [get]
func (controller *PlaceControllerImpl) PlaceFindByIdController(ctx place_presenter.PlacePresenterCTX) error {
	id := ctx.Param("id")

	place, err := controller.placeUseCase.FindPlaceByID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, place_presenter.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, place)
}

package place_controller_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	place_controller "github.com/risk-place-angola/backend-risk-place/api/rest/place/controllers"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	place_usecase "github.com/risk-place-angola/backend-risk-place/usecase/place"
	"github.com/stretchr/testify/assert"
)

type WsHandler struct {
	handler echo.HandlerFunc
}

func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	c := e.NewContext(r, w)

	forever := make(chan struct{})
	if err := h.handler(c); err != nil {
		log.Println(err)
	}

	<-forever
}

//nolint:errcheck
func TestPlaceFindAllWebSocket(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.Place{
		{
			ID:        "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			Latitude:  8.825248,
			Longitude: 13.263879,
		},
		{
			ID:        "50361691-6b99-8j2u-a8ca-db5d70912837",
			Latitude:  8.825248,
			Longitude: 13.263879,
		},
	}

	mockPlaceRepository := mocks.NewMockPlaceRepository(ctrl)
	mockPlaceRepository.EXPECT().FindAll().Return(data, nil)

	placeUseCase := place_usecase.NewPlaceUseCase(mockPlaceRepository)
	placeController := place_controller.NewPlaceClientManager(placeUseCase)

	go placeController.Start()

	h := WsHandler{
		handler: func(c echo.Context) error {
			return placeController.PlaceHandler(c)
		},
	}

	server := httptest.NewServer(http.HandlerFunc(h.ServeHTTP))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.Nil(t, err, err)

	// write
	err = ws.WriteMessage(websocket.TextMessage, []byte("ping"))
	assert.Nil(t, err, err)

	// read
	_, msg, err := ws.ReadMessage()

	json.Unmarshal(msg, &data)

	assert.Nil(t, err, err)
	assert.Equal(t, 2, len(data))

}

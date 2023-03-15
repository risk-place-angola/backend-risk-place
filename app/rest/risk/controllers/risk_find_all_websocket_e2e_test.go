package risk_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	risk_controller "github.com/risk-place-angola/backend-risk-place/app/rest/risk/controllers"
	"github.com/risk-place-angola/backend-risk-place/domain/entities"
	"github.com/risk-place-angola/backend-risk-place/domain/repository/mocks"
	risk_usecase "github.com/risk-place-angola/backend-risk-place/usecase/risk"
	"github.com/stretchr/testify/assert"
)

type WsHandler struct {
	handler echo.HandlerFunc
}

func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	c := e.NewContext(r, w)

	forever := make(chan struct{})
	h.handler(c)
	<-forever
}

func TestRiskFindAllWebSocket(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []*entities.Risk{
		{
			ID:          "93247691-5c64-4c1f-a8ca-db5d76640ca9",
			RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
			PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:        "Rangel rua da Lama",
			Latitude:    8.825248,
			Longitude:   13.263879,
			Description: "Risco de inundação",
		},
		{
			ID:          "50361691-6b99-8j2u-a8ca-db5d70912837",
			RiskTypeID:  "99bada49-09d0-4f13-b310-6f8633b38dfe",
			PlaceTypeID: "dd3aadda-9434-4dd7-aaad-035584b8f124",
			Name:        "Rangel rua da Lama",
			Latitude:    8.825248,
			Longitude:   13.263879,
			Description: "Risco de inundação",
		},
	}

	mockRiskRepository := mocks.NewMockRiskRepository(ctrl)
	mockRiskRepository.EXPECT().FindAll().Return(data, nil)

	riskUseCase := risk_usecase.NewRiskUseCase(mockRiskRepository)
	riskController := risk_controller.NewRiskClientManager(riskUseCase)

	go riskController.Start()

	h := WsHandler{
		handler: func(c echo.Context) error {
			return riskController.RiskHandler(c)
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

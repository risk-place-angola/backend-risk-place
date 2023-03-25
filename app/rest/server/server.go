package server

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/app/rest/dependency"
	"github.com/risk-place-angola/backend-risk-place/infra/db/drive/postgres"
)

type Server struct {
	Router *echo.Echo
}

func NewServer(router *echo.Echo) *Server {
	return &Server{
		Router: router,
	}
}

func (server *Server) Start() {

	server.Router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Angola!")
	})

	db, err := postgres.ConnectionPostgres()
	if err != nil {
		log.Panicln(err)
	}

	defer db.Close()

	dependency.Dependency(db, server.Router)

	server.Router.Logger.Fatal(server.Router.Start(":8000"))
}

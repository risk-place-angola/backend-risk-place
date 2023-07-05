package server

import (
	"github.com/jinzhu/gorm"
	"github.com/risk-place-angola/backend-risk-place/infra/db/drive/postgres"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/dependency"
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

	db, err := postgres.ConnectionPostgres()
	if err != nil {
		log.Panicln(err)
	}

	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	dependency.Dependency(db, server.Router)

	server.Router.Logger.Fatal(server.Router.Start(":8000"))
}

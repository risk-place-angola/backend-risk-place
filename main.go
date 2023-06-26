package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/risk-place-angola/backend-risk-place/infra/rest/server"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	Server := server.NewServer(echo.New())
	Server.Start()
}

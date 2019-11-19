package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/seungkyua/cookiemonster2/pkg/handler"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.PodHandler{}.SetHandler(e.Group("/api/v1/pod"))
	handler.ConfigHandler{}.SetHandler(e.Group("/api/v1/config"))
	handler.NodeHandler{}.SetHandler(e.Group("/api/v1/node"))

	e.Logger.Debug(e.Start(":10080"))
}

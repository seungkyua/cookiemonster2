package handler

import (
	"github.com/labstack/echo"
	"github.com/seungkyua/cookiemonster2/pkg/domain"
	"log"
	"net/http"
)

type NodeHandler struct{}

func (n NodeHandler) SetHandler(group *echo.Group) {
	group.GET("", n.List)
	group.POST("/start", n.NodeStart)
}

func (n NodeHandler) List(context echo.Context) error {
	log.Println("###########", randomInt(1))
	return context.JSONPretty(http.StatusOK, domain.Returnnamelist(), "    ")
}

func (n NodeHandler) NodeStart(context echo.Context) error {
	log.Println("###########", randomInt(1))
	return context.JSONPretty(http.StatusOK, domain.Reboot(), "    ")
}
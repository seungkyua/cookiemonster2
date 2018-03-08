package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/seungkyua/cookiemonster2/src/domain"
)

type ConfigHandler struct{}

/*
 * group.GET("", config.Get)
 */
func (h ConfigHandler) SetHandler(group *echo.Group) {
	group.GET("", h.Get)
}

// Get a config
func (h ConfigHandler) Get(context echo.Context) error {
	path := "config"
	c := domain.NewConfig()
	if err := c.ReadConfig(path); err != nil {
		fmt.Println(err)
	}

	return context.JSONPretty(http.StatusOK, c, "    ")
}

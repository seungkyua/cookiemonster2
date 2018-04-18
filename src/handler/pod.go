package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/seungkyua/cookiemonster2/src/domain"
	"golang.org/x/net/context"
)

type PodHandler struct{}

var m = &domain.PodManage{
	Started: false,
}

/*
 * group.GET("", pod.List)
 * group.GET("/:name", pod.Get)
 * group.DELETE("", pod.Delete)
 */

func (h PodHandler) SetHandler(group *echo.Group) {
	group.GET("", h.List)
	group.DELETE("", h.Delete)
	group.POST("", h.Stop)
}

// List running pods
func (h PodHandler) List(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World\n")
}

// Start a job to delete random pod
func (h PodHandler) Delete(c echo.Context) error {
	if m.Started {
		log.Println("Pods is already being munched, ignoring request\n")
		return nil
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		m.Ctx = ctx
		m.Cancel = cancel
		m.Started = true
	}

	err := m.Start(domain.GetConfig())
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func (h PodHandler) Stop(c echo.Context) error {
	if !m.Started {
		log.Println("Cookie Monster is not running, ignoring request\n")
		return nil
	}

	m.Stop(domain.GetConfig())

	return nil
}

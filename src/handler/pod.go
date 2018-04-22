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
 * group.GET("", PodHandler.List)
 * group.GET("/:name", PodHandler.Get)
 * group.POST("/start", PodHandler.Start)
 * group.POST("/stop", PodHandler.Stop)
 */
func (h PodHandler) SetHandler(group *echo.Group) {
	group.GET("", h.List)
	group.POST("/start", h.Start)
	group.POST("/stop", h.Stop)
}

// List running pods
func (h PodHandler) List(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World\n")
}

// Start a job to delete random pod
func (h PodHandler) Start(c echo.Context) error {
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
		log.Println(err)
		if m.Started {
			m.Stop(domain.GetConfig())
		}
	}

	return nil
}

func (h PodHandler) Stop(c echo.Context) error {
	if !m.Started {
		log.Println("Cookie Monster is currently munching, ignoring request\n")
		return nil
	}

	m.Stop(domain.GetConfig())
	log.Printf("Stop snacking.\n")

	return nil
}

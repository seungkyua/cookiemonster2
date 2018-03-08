package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/seungkyua/cookiemonster2/src/domain"
)

type PodHandler struct{}

/*
 * group.GET("", pod.List)
 * group.GET("/:name", pod.Get)
 * group.DELETE("", pod.Delete)
 */

func (h PodHandler) SetHandler(group *echo.Group) {
	group.GET("", h.List)
	group.DELETE("", h.Delete)
}

// List running pods
func (h PodHandler) List(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World\n")
}

// Start a job to delete random pod
func (h PodHandler) Delete(c echo.Context) error {
	DeletePod(domain.GetConfig())
	return nil
}

func DeletePod(config *domain.Config) {
	pod := &domain.VictimPod{}
	pod.GetVictimPod(config)
}

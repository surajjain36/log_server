package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatusResponse HTTP status
type StatusResponse struct {
	Name    string `json:"app"`
	Version string `json:"version"`
	Built   string `json:"built"`
	Status  bool   `json:"status"`
}

func (s *Service) ping(c *gin.Context) {
	status := true
	c.JSON(http.StatusOK, StatusResponse{
		Status:  status,
		Name:    s.AppName,
		Version: s.Version,
		Built:   s.BuildTime,
	})
}

func (s *Service) index(c *gin.Context) {
	c.String(http.StatusOK, "Nothing here")
}

func (s *Service) responseWriter(c *gin.Context, resp interface{}, code int) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(code, resp)
}

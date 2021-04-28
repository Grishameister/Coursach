package serverProxy

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Grishameister/Coursach/configs/config"
	"github.com/Grishameister/Coursach/internal/domain"
	"github.com/Grishameister/Coursach/internal/proxyHandlers"
	"github.com/Grishameister/Coursach/internal/ws"
)

type Server struct {
	config *config.ConfServer
	router *gin.Engine
}

func New(config *config.Config) *Server {
	statusChan := make(chan domain.StatusChannel, 256)
	handler := proxyHandlers.NewProxyHandler(
		http.Client{
			Timeout: time.Second * 10,
		}, statusChan)

	hub := ws.NewHub(statusChan)
	r := gin.Default()

	r.POST("/image", handler.HandleImages)
	r.GET("/image", handler.HandleImages)
	r.POST("/notification", handler.HandlerStatuses)
	r.GET("/ws", func(c *gin.Context) {
		hub.ServeWs(c.Writer, c.Request)
	})

	r.GET("/image/date", handler.HandleImages)
	r.GET("/image/last", handler.HandleImages)

	r.GET("/stat", handler.HandleStats)

	go hub.WriteMessages()

	return &Server{
		config: &config.Proxy.Server,
		router: r,
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf("%s:%s", s.config.Address, s.config.Port))
}

package serverProxy

import (
	"fmt"
	"github.com/Grishameister/Coursach/configs/config"
	"github.com/Grishameister/Coursach/internal/proxyHandlers"
	"github.com/Grishameister/Coursach/internal/tcpConnPool"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Server struct {
	config *config.ConfServer
	router *gin.Engine
}

func New(config *config.Config) *Server {

	pool, err := tcpConnPool.InitPool(10)

	if err != nil {
		log.Fatal("pool is not init")
	}

	handler := proxyHandlers.NewProxyHandler(
		http.Client{
		Timeout: time.Second * 10,
	}, pool)
	r := gin.Default()

	r.POST("/image", handler.HandleImages)
	r.GET("/image", handler.HandleImages)

	r.GET("/image/date", handler.HandleImages)
	r.GET("/image/last", handler.HandleImages)

	r.GET("/stat", handler.HandleStats)

	return &Server{
		config: &config.Proxy.Server,
		router: r,
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf("%s:%s", s.config.Address, s.config.Port))
}

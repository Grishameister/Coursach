package server

import (
	"fmt"
	"github.com/Grishamester/Coursach/configs/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.ConfServer
	router *gin.Engine
}

func New(config *config.Config) *Server {
	r := gin.Default()

	return &Server{
		config: &config.Web.Server,
		router: r,
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf("%s:%s", s.config.Address, s.config.Port))
}

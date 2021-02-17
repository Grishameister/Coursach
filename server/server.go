package server

import (
	"fmt"
	"github.com/Grishameister/Coursach/configs/config"
	"github.com/Grishameister/Coursach/internal/check"
	"github.com/Grishameister/Coursach/internal/database"
	"github.com/Grishameister/Coursach/internal/images/delivery"
	"github.com/Grishameister/Coursach/internal/images/repository"
	"github.com/Grishameister/Coursach/internal/images/usecase"
	"github.com/Grishameister/Coursach/internal/poolDb"
	"github.com/Grishameister/Coursach/internal/queue"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.ConfServer
	router *gin.Engine
}



func New(config *config.Config, db database.DBInterface, queue *queue.Queue) *Server {
	rep := repository.NewRepo(db)
	uc := usecase.NewUsecase(rep)
	ch := make(chan interface{})
	in := make (chan []byte, 100)
	handler := delivery.NewHandler(queue, ch, in)

	pool := poolDb.NewPool(in, uc)
	pool.Run(8)

	go check.Check(ch)

	r := gin.Default()

	r.POST("/image", handler.SaveFrameMiddleWare(), handler.ToQueue)
	r.GET("/image", handler.FromQueue)

	r.GET("/image/date")
	r.GET("/image/last")

	return &Server{
		config: &config.Web.Server,
		router: r,
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf("%s:%s", s.config.Address, s.config.Port))
}

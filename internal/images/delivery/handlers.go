package delivery

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/Grishameister/Coursach/internal/images"
	"github.com/Grishameister/Coursach/internal/queue"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	q *queue.Queue
	uc images.IUsecase
}

func NewHandler(q *queue.Queue, uc images.IUsecase) *Handler {
	return &Handler{
		q:  q,
		uc: uc,
	}
}

func (h *Handler) SaveFrameMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("bad form"))
			return
		}

		src, err := file.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("error"))
			return
		}
		defer src.Close()

		reader := bufio.NewReader(src)

		buffer, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Println("Cant read")
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("error"))
			return
		}
		c.Set("buffer", buffer)
		c.Next()

		var zipbuf bytes.Buffer
		gz := gzip.NewWriter(&zipbuf)

		if _, err := gz.Write(buffer); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if err := gz.Close(); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := h.uc.PostFrame(c, zipbuf.Bytes()); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) ToQueue(c *gin.Context) {
	buffer, ok := c.Get("buffer")
	if !ok {
		log.Println("Buffer doesnt work")
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	if err := h.q.Push(buffer.([]byte)); err != nil {
		log.Println(err, "Push doesnt work")
		c.AbortWithStatus(http.StatusBadGateway)
	}
	c.Status(200)
}

func (h *Handler) FromQueue(c *gin.Context) {
	data, err := h.q.Pop()
	if err != nil {
		log.Println("CANT POP FROM QUEUE")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	_, err = c.Writer.Write(data)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
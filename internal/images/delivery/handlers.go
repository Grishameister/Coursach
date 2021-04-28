package delivery

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/Grishameister/Coursach/internal/check"
	"github.com/Grishameister/Coursach/internal/images"
	"github.com/Grishameister/Coursach/internal/queue"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Handler struct {
	q      *queue.Queue
	notify chan interface{}
	toPool chan []byte
	uc     images.IUsecase
}

func NewHandler(q *queue.Queue, ch chan interface{}, toPool chan []byte, uc images.IUsecase) *Handler {
	return &Handler{
		q:      q,
		notify: ch,
		toPool: toPool,
		uc:     uc,
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
		h.toPool <- buffer
	}
}

func (h *Handler) ToQueue(c *gin.Context) {
	buffer, ok := c.Get("buffer")
	if !ok {
		log.Println("Buffer doesnt work")
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	if !check.PassToQueue {
		log.Println("Wait for ")
		c.Status(http.StatusOK)
		return
	}
	if err := h.q.Push(buffer.([]byte)); err != nil {
		log.Println(err, "Push doesnt work")
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	c.Status(200)
	log.Println("SIZE OF QUEUE IS ", h.q.Size())
}

func WriteRequestData(c *gin.Context, data *[]byte) {
	if data == nil {
		return
	}
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Length", strconv.Itoa(len(*data)))
	_, err := c.Writer.Write(*data)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) FromQueue(c *gin.Context) {
	data, err := h.q.Pop()
	h.notify <- struct{}{}
	if err != nil {
		log.Println("CANT POP FROM QUEUE")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	WriteRequestData(c, &data)
	log.Println("SIZE OF QUEUE IS ", h.q.Size())
}

func WriteGzipData(c *gin.Context, data []byte) {
	if data == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	bufData := bytes.NewReader(data)
	gz, err := gzip.NewReader(bufData)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	size, err := io.Copy(c.Writer, gz)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer gz.Close()
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Length", strconv.FormatInt(size, 10))
	c.Status(http.StatusOK)
}

func (h *Handler) GetLastFrame(c *gin.Context) {
	WriteGzipData(c, h.uc.GetLastFrame())
}

func (h *Handler) GetFrameByDate(c *gin.Context) {
	reqDate, err := url.QueryUnescape(c.Query("date"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02 15:04:01", reqDate)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	WriteGzipData(c, h.uc.GetFrameByDate(date))
}

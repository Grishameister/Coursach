package proxyHandlers

import (
	"github.com/Grishameister/Coursach/configs/config"
	"github.com/Grishameister/Coursach/internal/domain"
	"github.com/Grishameister/Coursach/internal/tcpConnPool"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ProxyHandler struct {
	client http.Client
	p *tcpConnPool.TcpPool
}

func NewProxyHandler(client http.Client, p *tcpConnPool.TcpPool) *ProxyHandler {
	return &ProxyHandler{
		client: client,
		p: p,
	}
}

func (h *ProxyHandler) HandleImages(c *gin.Context) {
	req, err := http.NewRequest(c.Request.Method, "http://" + config.Conf.Web.Server.Address + ":"+ config.Conf.Web.Server.Port + c.Request.RequestURI, c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		return
	}

	for key, values := range c.Request.Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	req.Host = config.Conf.Web.Server.Address + ":" + config.Conf.Web.Server.Port

	resp, err := h.client.Do(req)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	for key, values := range resp.Header {
		for _, v := range values {
			c.Writer.Header().Add(key, v)
		}
	}

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c.Status(resp.StatusCode)
}

func (h *ProxyHandler) HandleStats(c *gin.Context) {
	reqDate, err := url.QueryUnescape(c.Query("date"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02 15:04", reqDate)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	s := domain.GetStat{
		Date: date,
	}

	b, err := msgpack.Marshal(&s)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	conn, err := h.p.Pop()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = conn.Write(b)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resp := make([]byte, 1024)

	n, err := conn.Read(resp)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := domain.StatFromServer{}

	if err := msgpack.Unmarshal(resp[0:n], &response); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.p.Push(conn); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, response)
}

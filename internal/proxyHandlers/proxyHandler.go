package proxyHandlers

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/Grishameister/Coursach/configs/config"
	"github.com/Grishameister/Coursach/internal/domain"
	"github.com/Grishameister/Coursach/internal/tcpConnPool"
)

type ProxyHandler struct {
	client     http.Client
	p          *tcpConnPool.TcpPool
	ch         chan domain.StatusChannel
	LastStatus domain.Statuses
	mu         *sync.RWMutex
}

func NewProxyHandler(client http.Client, ch chan domain.StatusChannel) *ProxyHandler {
	return &ProxyHandler{
		client: client,
		ch:     ch,
		LastStatus: domain.Statuses{
			domain.StatusOK: struct{}{},
		},
		mu: &sync.RWMutex{},
	}
}

func (h *ProxyHandler) HandleImages(c *gin.Context) {
	req, err := http.NewRequest(c.Request.Method, "http://"+config.Conf.Web.Server.Address+":"+config.Conf.Web.Server.Port+c.Request.RequestURI, c.Request.Body)
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

	s := domain.DataRequest{
		Date:      reqDate,
		TypeRoute: domain.ReadStatData,
	}

	b, err := msgpack.Marshal(s)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	conn, err := net.Dial("tcp", config.Conf.Stats.Server.Address+":"+config.Conf.Stats.Server.Port)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

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

	response := domain.DataResponse{}
	if err := msgpack.Unmarshal(resp[0:n], &response); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProxyHandler) HandlerStatuses(c *gin.Context) {
	var req domain.StatusRequest
	if err := msgpack.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	statuses := strings.Split(req.Statuses, ",")

	h.mu.Lock()
	if different := h.checkStatuses(statuses); different {
		h.LastStatus = newStatus(statuses)
		h.ch <- domain.StatusChannel{
			Date:     req.Date,
			Statuses: statuses,
		}
	}
	h.mu.Unlock()

	c.Status(http.StatusOK)
}

func (h *ProxyHandler) checkStatuses(statuses []string) bool {
	different := false
	for _, st := range statuses {
		if _, ok := h.LastStatus[domain.Status(st)]; !ok {
			different = true
			break
		}
	}
	return different
}

func newStatus(statuses []string) domain.Statuses {
	result := domain.Statuses{}
	for _, st := range statuses {
		result[domain.Status(st)] = struct{}{}
	}
	return result
}

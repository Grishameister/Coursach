package proxyHandlers

import (
	"github.com/Grishameister/Coursach/configs/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

type ProxyHandler struct {
	client http.Client
}

func NewProxyHandler(client http.Client) *ProxyHandler {
	return &ProxyHandler{
		client: client,
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

	req.Host = "localhost:8008"

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

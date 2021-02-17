package proxyHandlers

import (
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
	resp, err := h.client.Do(c.Request)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	for key, values := range resp.Header {
		for _, v := range values {
			c.Header(key, v)
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

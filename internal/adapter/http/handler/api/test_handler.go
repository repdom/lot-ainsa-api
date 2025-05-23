package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TestHandler struct {
}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h *TestHandler) TestHandler(c *gin.Context) {
	log.Printf("Test handler called")
	c.Status(http.StatusOK)
}

package api

import (
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type FinancingsActiveHandler struct {
	services port.FinancingsActionService
}

func (handler *FinancingsActiveHandler) HandleRequest(c *gin.Context) {

	var loan model.RequestLoan

	financingIdQuery := c.Query("financingId")

	financingId, err := strconv.Atoi(financingIdQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "financingId not valid"})
		return
	}

	if err := c.ShouldBindJSON(&loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt := c.Request.Header.Get("Authorization")
	user := c.Request.Header.Get("x-user")
	lang := c.Request.Header.Get("x-language")

	log.Println(jwt)
	log.Println(user)
	log.Println(lang)

	err = handler.services.Activation(loan, financingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func NewFinancingHandler(services port.FinancingsActionService) *FinancingsActiveHandler {
	return &FinancingsActiveHandler{
		services: services,
	}
}

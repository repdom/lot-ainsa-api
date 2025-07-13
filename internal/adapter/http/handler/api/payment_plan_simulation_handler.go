package api

import (
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type PaymentPlanSimulationHandler struct {
	service port.ApiService
}

func NewPaymentSimulationHandler(service port.ApiService) *PaymentPlanSimulationHandler {
	return &PaymentPlanSimulationHandler{service: service}
}

func (handler *PaymentPlanSimulationHandler) HandleRequest(c *gin.Context) {

	var loan model.RequestLoan

	jwt := c.Request.Header.Get("Authorization")
	user := c.Request.Header.Get("x-user")
	lang := c.Request.Header.Get("x-language")

	log.Println(jwt)
	log.Println(user)
	log.Println(lang)

	if err := c.ShouldBindJSON(&loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, err := handler.service.GenerateSimulation(loan)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, body)

}

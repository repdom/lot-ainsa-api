package api

import (
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type CustomerOnboardingHandler struct {
	service port.CustomerOnboardingService
}

func NewCustomerOnboardingHandler(service port.CustomerOnboardingService) *CustomerOnboardingHandler {
	return &CustomerOnboardingHandler{service: service}
}

func (handler *CustomerOnboardingHandler) CreateCustomer(c *gin.Context) {
	var customer model.RequestCustomerOnboarding

	jwt := c.Request.Header.Get("Authorization")
	user := c.Request.Header.Get("x-user")
	lang := c.Request.Header.Get("x-language")

	log.Println(jwt)
	log.Println(user)
	log.Println(lang)

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	badRequest, internalServer := handler.service.CreateCustomer(jwt, user, lang, customer)

	if internalServer != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalServer.Error()})
		return
	}

	if badRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": badRequest.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

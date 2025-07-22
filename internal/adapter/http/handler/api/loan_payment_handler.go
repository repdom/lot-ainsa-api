package api

import (
	"be-lotsanmateo-api/internal/domain/port"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type LoanPaymentHandler struct {
	loanPayment port.LoanPaymentService
}

func NewLoanPaymentHandler(loanPayment port.LoanPaymentService) *LoanPaymentHandler {
	return &LoanPaymentHandler{
		loanPayment: loanPayment,
	}
}

func (handler *LoanPaymentHandler) HandleRequest(c *gin.Context) {

	financingIdQuery := c.Query("financingId")

	financingId, err := strconv.Atoi(financingIdQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "financingId not valid"})
		return
	}

	share := c.Query("share")

	shareAmount, err := strconv.ParseFloat(share, 64)

	jwt := c.Request.Header.Get("Authorization")
	user := c.Request.Header.Get("x-user")
	lang := c.Request.Header.Get("x-language")

	log.Println(jwt)
	log.Println(user)
	log.Println(lang)

	resp, err := handler.loanPayment.CalculateLoanPayment(financingId, shareAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)

}

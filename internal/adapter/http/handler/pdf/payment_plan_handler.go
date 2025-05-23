package pdf

import (
	"be-lotsanmateo-api/internal/domain/port"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PaymentPlanDocumentHandler struct {
	service port.ReportService
}

func NewPaymentPlanDocumentHandler(service port.ReportService) *PaymentPlanDocumentHandler {
	return &PaymentPlanDocumentHandler{service: service}
}

func (handler *PaymentPlanDocumentHandler) GeneratePDF(c *gin.Context) {

	lotIdRequest := c.Query("lotId")
	customerIdRequest := c.Query("customerId")
	view := c.Query("view")

	lotId, err := strconv.Atoi(lotIdRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lotId not valid"})
		return
	}

	customerId, err := strconv.Atoi(customerIdRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lotId not valid"})
		return
	}

	boolValue, err := strconv.ParseBool(view)
	val := ""
	if boolValue {
		val += "attachment;"
	} else {
		val += "inline;"
	}

	val += "filename=plan_de_pago.pdf"

	pdfData, err := handler.service.GenerateReport(lotId, customerId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar PDF"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", val)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

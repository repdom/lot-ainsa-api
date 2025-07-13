package pdf

import (
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
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

	financingIdQuery := c.Query("financingId")
	view := c.Query("view")

	financingId, err := strconv.Atoi(financingIdQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "financingId not valid"})
		return
	}

	boolValue, err := strconv.ParseBool(view)
	val := ""
	if !boolValue {
		val += "attachment;"
	} else {
		val += "inline;"
	}

	val += "filename=plan_de_pago.pdf"

	pdfData, err := handler.service.GenerateReport(financingId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar PDF: %s", err.Error())})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", val)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

package pdf

import (
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentPlanSimulationDocumentHandler struct {
	service port.ReportService
}

func NewPaymentPlanSimulationDocumentHandler(service port.ReportService) *PaymentPlanSimulationDocumentHandler {
	return &PaymentPlanSimulationDocumentHandler{service: service}
}

func (handler *PaymentPlanSimulationDocumentHandler) GeneratePDF(c *gin.Context) {

	lotIdRequest := c.Query("lotId")
	view := c.Query("view")

	jwt := c.Request.Header.Get("Authorization")
	user := c.Request.Header.Get("x-user")
	lang := c.Request.Header.Get("x-language")

	lotId, err := strconv.Atoi(lotIdRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lotId not valid"})
		return
	}

	boolValue, err := strconv.ParseBool(view)
	val := ""
	if !boolValue {
		val += "attachment;"
	} else {
		val += "inline;"
	}

	pdfData, name, err := handler.service.GenerateReport(jwt, user, lang, lotId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar PDF"})
		return
	}

	val += fmt.Sprintf("filename=%s_plan_de_pago.pdf", *name)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", val)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

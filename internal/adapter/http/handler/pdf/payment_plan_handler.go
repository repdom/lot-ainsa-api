package pdf

import (
	"be-lotsanmateo-api/internal/adapter/http/handler/pdf/utility"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentPlanDocumentHandler struct {
	service port.ReportService
}

func NewPaymentPlanDocumentHandler(service port.ReportService) *PaymentPlanDocumentHandler {
	return &PaymentPlanDocumentHandler{service: service}
}

func (handler *PaymentPlanDocumentHandler) GeneratePDF(c *gin.Context) {

	financingId, ok := utility.ParseQueryInt(c, "financingId")
	view, _ := utility.ParseQueryBool(c, "view")
	if !ok {
		return
	}

	jwt, user, lang, _ := utility.ExtractHeaders(c)

	pdfData, name, err := handler.service.GenerateReport(jwt, user, lang, financingId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar PDF: %s", err.Error())})
		return
	}

	disposition := utility.BuildDispositionBool(view, fmt.Sprintf("filename=%s_plan_de_pago.pdf", *name))

	utility.Response(c, disposition, pdfData)
}

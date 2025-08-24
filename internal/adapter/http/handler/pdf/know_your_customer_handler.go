package pdf

import (
	"be-lotsanmateo-api/internal/adapter/http/handler/pdf/utility"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KnowYourCustomerHandler struct {
	service port.KnowYourCustomerService
}

func NewKnowYourCustomerHandler(service port.ReportService) *KnowYourCustomerHandler {
	return &KnowYourCustomerHandler{service: service}
}

func (handler *KnowYourCustomerHandler) GeneratePDF(c *gin.Context) {

	customerId, ok := utility.ParseQueryInt(c, "customerId")
	view, _ := utility.ParseQueryBool(c, "view")
	if !ok {
		return
	}

	jwt, user, lang, _ := utility.ExtractHeaders(c)

	pdfData, name, err := handler.service.GenerateReport(jwt, user, lang, customerId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar PDF: %s", err.Error())})
		return
	}

	disposition := utility.BuildDispositionBool(view, fmt.Sprintf("filename=%s.pdf", *name))

	utility.Response(c, disposition, pdfData)
}

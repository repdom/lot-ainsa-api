package pdf

import (
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InvoiceDocumentHandler struct {
	service port.InvoiceService
}

func NewInvoiceDocumentHandler(service port.InvoiceService) *InvoiceDocumentHandler {
	return &InvoiceDocumentHandler{service: service}
}

func (handler *InvoiceDocumentHandler) GeneratePayment(c *gin.Context) {
	jwt, user, lang, name := extractHeaders(c)
	reservationId, ok := parseQueryInt(c, "paymentId")
	if !ok {
		return
	}

	var pdfData []byte
	pdfData, clientName, err := handler.service.InvoicePayment(jwt, user, lang, name, reservationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar la factura de pago PDF: %s", err.Error())})
		return
	}

	disposition := buildDisposition(c.Query("view"), fmt.Sprintf("%s_pago.pdf", *clientName))

	response(c, disposition, pdfData)
}

func response(c *gin.Context, disposition string, pdfData []byte) {
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", disposition)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

func (handler *InvoiceDocumentHandler) GenerateDownPayment(c *gin.Context) {
	jwt, user, lang, name := extractHeaders(c)
	reservationId, ok := parseQueryInt(c, "downPaymentId")
	if !ok {
		return
	}

	var pdfData []byte
	pdfData, clientName, err := handler.service.InvoiceDownPayment(jwt, user, lang, name, reservationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar la factura de la prima PDF: %s", err.Error())})
		return
	}

	disposition := buildDisposition(c.Query("view"), fmt.Sprintf("%s_prima.pdf", *clientName))

	response(c, disposition, pdfData)
}

func (handler *InvoiceDocumentHandler) GenerateReservation(c *gin.Context) {
	jwt, user, lang, name := extractHeaders(c)
	reservationId, ok := parseQueryInt(c, "reservationId")
	if !ok {
		return
	}

	var pdfData []byte
	pdfData, clientName, err := handler.service.InvoiceReservation(jwt, user, lang, name, reservationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar la factura de la reserva PDF: %s", err.Error())})
		return
	}

	disposition := buildDisposition(c.Query("view"), fmt.Sprintf("%s_reserva.pdf", *clientName))

	response(c, disposition, pdfData)
}

func extractHeaders(c *gin.Context) (jwt, user, lang, name string) {
	jwt = c.GetHeader("Authorization")
	user = c.GetHeader("x-user")
	lang = c.GetHeader("x-language")
	name = c.GetString("name")
	return
}

func parseQueryInt(c *gin.Context, key string) (int, bool) {
	valueStr := c.Query(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s not valid", key)})
		return 0, false
	}
	return value, true
}

func buildDisposition(view string, filename string) string {
	isInline, err := strconv.ParseBool(view)
	dispositionType := "attachment;"
	if err == nil && isInline {
		dispositionType = "inline;"
	}
	return fmt.Sprintf("%sfilename=%s", dispositionType, filename)
}

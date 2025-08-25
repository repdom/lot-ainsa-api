package pdf

import (
	"be-lotsanmateo-api/internal/adapter/http/handler/pdf/utility"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvoiceDocumentHandler struct {
	service port.InvoiceService
}

func NewInvoiceDocumentHandler(service port.InvoiceService) *InvoiceDocumentHandler {
	return &InvoiceDocumentHandler{service: service}
}

func (handler *InvoiceDocumentHandler) GeneratePayment(c *gin.Context) {
	jwt, user, lang, name := utility.ExtractHeaders(c)
	reservationId, ok := utility.ParseQueryInt(c, "paymentId")
	if !ok {
		return
	}

	var pdfData []byte
	pdfData, clientName, err := handler.service.InvoicePayment(jwt, user, lang, name, reservationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar la factura de pago PDF: %s", err.Error())})
		return
	}

	disposition := utility.BuildDisposition(c.Query("view"), fmt.Sprintf("%s_pago.pdf", *clientName))

	utility.Response(c, disposition, pdfData)
}

func (handler *InvoiceDocumentHandler) GenerateDownPayment(c *gin.Context) {
	jwt, user, lang, name := utility.ExtractHeaders(c)
	reservationId, ok := utility.ParseQueryInt(c, "downPaymentId")
	if !ok {
		return
	}

	var pdfData []byte
	pdfData, clientName, err := handler.service.InvoiceDownPayment(jwt, user, lang, name, reservationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar la factura de la prima PDF: %s", err.Error())})
		return
	}

	disposition := utility.BuildDisposition(c.Query("view"), fmt.Sprintf("%s_prima.pdf", *clientName))

	utility.Response(c, disposition, pdfData)
}

func (handler *InvoiceDocumentHandler) GenerateReservation(c *gin.Context) {
	jwt, user, lang, name := utility.ExtractHeaders(c)
	reservationId, ok := utility.ParseQueryInt(c, "reservationId")
	if !ok {
		return
	}

	var pdfData []byte
	pdfData, clientName, err := handler.service.InvoiceReservation(jwt, user, lang, name, reservationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar la factura de la reserva PDF: %s", err.Error())})
		return
	}

	disposition := utility.BuildDisposition(c.Query("view"), fmt.Sprintf("%s_reserva.pdf", *clientName))

	utility.Response(c, disposition, pdfData)
}

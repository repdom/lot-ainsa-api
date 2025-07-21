package pdf

import (
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type PagareHandler struct {
	pagare pdf.GeneratePagarePDF
}

func NewPagareHandler() *PagareHandler {
	return &PagareHandler{
		pagare: pdf.NewPagarePDF(),
	}
}

func (handler *PagareHandler) GeneratePDF(c *gin.Context) {

	financingIdQuery := c.Query("financingId")
	view := c.Query("view")

	financingId, err := strconv.Atoi(financingIdQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "financingId not valid"})
		return
	}
	log.Println(financingId)

	boolValue, err := strconv.ParseBool(view)
	val := ""
	if !boolValue {
		val += "attachment;"
	} else {
		val += "inline;"
	}

	val += "filename=plan_de_pago.pdf"

	data := pdf.PagareData{}

	data.Area = "200"
	data.Address = "Calle 123 # 123"

	pdfData, err := handler.pagare.GenerateReport(data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar PDF: %s", err.Error())})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", val)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

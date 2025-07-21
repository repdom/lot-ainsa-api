package pdf

import (
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type PagareHandler struct {
	pagare port.PagareService
}

func NewPagareHandler(pagare port.PagareService) *PagareHandler {
	return &PagareHandler{
		pagare: pagare,
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

	val += "filename=pagare.pdf"

	pdfData, err := handler.pagare.GenerateReport(financingId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar PDF: %s", err.Error())})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", val)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

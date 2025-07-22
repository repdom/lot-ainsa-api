package pdf

import (
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type PromissoryNoteHandler struct {
	promissoryNote port.PromissoryNoteService
}

func NewPromissoryNoteHandler(promissoryNote port.PromissoryNoteService) *PromissoryNoteHandler {
	return &PromissoryNoteHandler{
		promissoryNote: promissoryNote,
	}
}

func (handler *PromissoryNoteHandler) GeneratePDF(c *gin.Context) {

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

	pdfData, err := handler.promissoryNote.GenerateReport(financingId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar PDF: %s", err.Error())})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", val)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

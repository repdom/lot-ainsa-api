package pdf

import (
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	jwt := c.Request.Header.Get("Authorization")
	user := c.Request.Header.Get("x-user")
	lang := c.Request.Header.Get("x-language")

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

	pdfData, name, err := handler.promissoryNote.GenerateReport(jwt, user, lang, financingId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al generar PDF: %s", err.Error())})
		return
	}

	val += fmt.Sprintf("filename=%s_plan_de_pago.pdf", *name)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", val)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

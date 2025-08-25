package utility

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ExtractHeaders(c *gin.Context) (jwt, user, lang, name string) {
	jwt = c.GetHeader("Authorization")
	user = c.GetHeader("x-user")
	lang = c.GetHeader("x-language")
	name = c.GetString("name")
	return
}

func ParseQueryInt(c *gin.Context, key string) (int, bool) {
	valueStr := c.Query(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s not valid", key)})
		return 0, false
	}
	return value, true
}

func ParseQueryBool(c *gin.Context, key string) (bool, bool) {
	valueStr := c.Query(key)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, true
	}
	return value, true
}

func BuildDisposition(view string, filename string) string {
	isInline, err := strconv.ParseBool(view)
	dispositionType := "attachment;"
	if err == nil && isInline {
		dispositionType = "inline;"
	}
	return fmt.Sprintf("%sfilename=%s", dispositionType, filename)
}

func BuildDispositionBool(view bool, filename string) string {
	dispositionType := "attachment;"
	if view {
		dispositionType = "inline;"
	}
	return fmt.Sprintf("%sfilename=%s", dispositionType, filename)
}

func Response(c *gin.Context, disposition string, pdfData []byte) {
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", disposition)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

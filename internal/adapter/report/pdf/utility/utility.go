package utility

import (
	"bytes"
	"html/template"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// 1. Cargar plantilla desde archivo
func LoadTemplate(path string) (*template.Template, error) {
	tpl, err := template.ParseFiles(path)
	if err != nil {
		log.Println("Error al cargar plantilla:", err)
		return nil, err
	}
	return tpl, nil
}

func ExecuteTemplateToHTML(tpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		log.Println("Error al ejecutar plantilla:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

type Margin struct {
	MarginTop    uint
	MarginBottom uint
	MarginLeft   uint
	MarginRight  uint
}

func NewMarginDefault() *Margin {
	return &Margin{
		MarginTop:    0,
		MarginBottom: 0,
		MarginLeft:   0,
		MarginRight:  0,
	}
}

func NewMargin(top, bottom, left, right uint) *Margin {
	return &Margin{
		MarginTop:    top,
		MarginBottom: bottom,
		MarginLeft:   left,
		MarginRight:  right,
	}
}

func GeneratePDFFromHTML(html []byte, margin *Margin) ([]byte, error) {
	diff, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println("Error al crear PDF generator:", err)
		return nil, err
	}

	// Configuración básica
	if margin != nil {
		diff.MarginTop.Set(margin.MarginTop)
		diff.MarginBottom.Set(margin.MarginBottom)
		diff.MarginLeft.Set(margin.MarginLeft)
		diff.MarginRight.Set(margin.MarginRight)
	}

	diff.Dpi.Set(720)

	diff.NoCollate.Set(false)
	diff.PageSize.Set(wkhtmltopdf.PageSizeLetter)

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(html))
	page.EnableLocalFileAccess.Set(true)
	page.FooterLine.Set(true)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(12)

	diff.AddPage(page)

	if err := diff.Create(); err != nil {
		log.Println("Error al generar PDF:", err)
		return nil, err
	}

	return diff.Bytes(), nil
}

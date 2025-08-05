package pdf

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"log"
)

type GeneratePagarePDF struct{}

type PagareData struct {
	// {{codigo de pagare}}
	NoteCode string

	// {{numero de lote}}
	LotNumber string

	Area string

	// {{poligono}}
	Block string

	// {{Nombre del cliente}}
	ClientName string

	// {{edad del cliente}}
	ClientAge int

	// {{profesion}}
	Profession string

	// {{direccion}}
	Address string

	// {{documento de identidad}}
	IdentityDocument string

	DateNow string

	// {{monto costo del lote financiado}}
	LotCost string

	// {{monto costo del lote en letras}}
	LotCostInWords string

	// {{monto de la prima dada}}
	DownPayment string

	// {{monto de la prima dada en letras}}
	DownPaymentInWords string

	// {{porcentaje de la prima dada}}
	DownPaymentPercentage string

	// {{monto a financiar}}
	AmountToFinance string

	// {{monto a financiar en letras}}
	AmountToFinanceInWords string

	// {{plazo en años}}
	TermInYears int

	// {{cantidad de cuotas}}
	InstallmentCount int

	// {{tasa de interes aplicada}}
	InterestRate      string
	InterestOnArrears string

	// {{monto de la cuota}}
	InstallmentAmount string

	// {{fecha del primer pago de la cuota}}
	FirstPaymentDate string
}

func (p GeneratePagarePDF) GenerateReport(data PagareData) ([]byte, error) {
	// Cargar la plantilla
	tpl, err := template.ParseFiles("docs/pagare.html")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Ejecutar plantilla con datos
	var htmlBuffer bytes.Buffer
	if err := tpl.Execute(&htmlBuffer, data); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Generar PDF
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	pdfg.MarginTop.Set(20)
	pdfg.MarginBottom.Set(20)
	pdfg.MarginLeft.Set(20)
	pdfg.MarginRight.Set(20)

	pdfg.Dpi.Set(720)
	pdfg.NoCollate.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeLetter)

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes()))
	page.EnableLocalFileAccess.Set(true)
	page.FooterLine.Set(true)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(24)

	pdfg.AddPage(page)

	if err := pdfg.Create(); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return pdfg.Bytes(), nil
}

func NewGeneratePagarePDF() GeneratePagarePDF {
	return GeneratePagarePDF{}
}

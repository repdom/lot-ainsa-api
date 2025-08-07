package pdf

import (
	"bytes"
	"html/template"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type InvoicePagarePDF struct{}

type ReciboPagoData struct {
	CantidadPagadaEnNumeros string
	LugarYFecha             string
	RecibiDe                string
	LaCantidadDe            string
	Direccion               string
	Dui                     string
	Nit                     string
	ConceptoDe              string
	LoteN                   string
	UbicacionDePoligono     string
	Efectivo                string
	TransferenciaID         string
	FechaTransferencia      string
	BancoT                  string
	BancoR                  string
	Ctan                    string
	NumeroDeposito          string
	DepositoBanco           string
	RecibidoPor             string
	Firma                   string
	ValorDelLote            string
	SaldoAnterior           string
	AbonoACapital           string
	InteresesNormales       string
	InteresesMoratorios     string
	SaldoActual             string
}

func (p InvoicePagarePDF) GenerateReport(data ReciboPagoData) ([]byte, error) {
	// Cargar la plantilla
	tpl, err := template.ParseFiles("docs/factura.gohtml")
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

	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)

	pdfg.Dpi.Set(720)
	pdfg.NoCollate.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeLetter)

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes()))
	page.EnableLocalFileAccess.Set(true)
	page.FooterLine.Set(true)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(12)

	pdfg.AddPage(page)

	if err := pdfg.Create(); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return pdfg.Bytes(), nil
}

func NewInvoicePagarePDF() *InvoicePagarePDF {
	return &InvoicePagarePDF{}
}

package pdf

import (
	"be-lotsanmateo-api/internal/adapter/report/pdf/utility"
)

type InvoicePagarePDF struct{}

type ReciboPagoData struct {
	EmpresaNombre           string
	EmpresaDireccion        string
	EmpresaTelefono         string
	EmpresaEmail            string
	FechaEmision            string
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

	tpl, err := utility.LoadTemplate("docs/factura.gohtml")
	if err != nil {
		return nil, err
	}

	html, err := utility.ExecuteTemplateToHTML(tpl, data)
	if err != nil {
		return nil, err
	}

	return utility.GeneratePDFFromHTML(html, utility.NewMarginDefault())
}

func NewInvoicePagarePDF() *InvoicePagarePDF {
	return &InvoicePagarePDF{}
}

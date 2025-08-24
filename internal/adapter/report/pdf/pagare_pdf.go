package pdf

import (
	"be-lotsanmateo-api/internal/adapter/report/pdf/utility"
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

	FinancingAmountPercentage string

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
	tpl, err := utility.LoadTemplate("docs/pagare.gohtml")
	if err != nil {
		return nil, err
	}

	html, err := utility.ExecuteTemplateToHTML(tpl, data)
	if err != nil {
		return nil, err
	}

	return utility.GeneratePDFFromHTML(html, utility.NewMarginDefault())
}

func NewGeneratePagarePDF() GeneratePagarePDF {
	return GeneratePagarePDF{}
}

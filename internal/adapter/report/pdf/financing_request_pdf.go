package pdf

import "be-lotsanmateo-api/internal/adapter/report/pdf/utility"

type FinancingRequestPDF struct {
}

func NewFinancingRequestPDF() *FinancingRequestPDF {
	return &FinancingRequestPDF{}
}

type SolicitudFinanciamiento struct{}

func (p FinancingRequestPDF) GenerateReport(data SolicitudFinanciamiento) ([]byte, error) {

	tpl, err := utility.LoadTemplate("docs/conoceTuCliente.gohtml")
	if err != nil {
		return nil, err
	}

	html, err := utility.ExecuteTemplateToHTML(tpl, data)
	if err != nil {
		return nil, err
	}

	return utility.GeneratePDFFromHTML(html, utility.NewMarginDefault())
}

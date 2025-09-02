package pdf

import "be-lotsanmateo-api/internal/adapter/report/pdf/utility"

type FinancingRequestPDF struct {
}

func NewFinancingRequestPDF() *FinancingRequestPDF {
	return &FinancingRequestPDF{}
}

type SolicitudFinanciamiento struct {
	Fecha                       string
	Ciudad                      string
	Departamento                string
	NombreUrbanizacion          string
	NombreDestinatario          string
	PorcentajeFinanciamiento    string
	PorcentajeFinanciamientoTxt string
	MontoFinanciamiento         string
	Plazo                       string
	NombreCompleto              string
	NumeroDUI                   string
	DireccionCompleta           string
	ActividadEconomica          string
	NumeroLote                  string
	Poligono                    string
	PrecioLote                  string
	AbonoEfectivo               string
	Telefono                    string
	CorreoElectronico           string
}

func (p FinancingRequestPDF) GenerateReport(data SolicitudFinanciamiento) ([]byte, error) {

	tpl, err := utility.LoadTemplate("docs/solicitud.gohtml")
	if err != nil {
		return nil, err
	}

	html, err := utility.ExecuteTemplateToHTML(tpl, data)
	if err != nil {
		return nil, err
	}

	return utility.GeneratePDFFromHTML(html, utility.NewMarginDefault())
}

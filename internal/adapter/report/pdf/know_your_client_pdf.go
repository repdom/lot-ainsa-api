package pdf

import "be-lotsanmateo-api/internal/adapter/report/pdf/utility"

type ConoceTuCliente struct {
	NombreCompleto         string
	TipoDocumento          string
	NumeroDocumento        string
	FechaNacimiento        string
	Genero                 string
	EstadoCivil            string
	Nacionalidad           string
	DireccionResidencia    string
	Ciudad                 string
	CodigoPostal           string
	TelefonoContacto       string
	CorreoElectronico      string
	OcupacionProfesion     string
	NombreEmpresa          string
	CargoPuesto            string
	FuenteIngresos         string
	RangoIngresosMensuales string
	PropositoRelacion      string
	EsPEP                  string
	DetallesPEP            string
	Firma                  string
}

type KnowYourClientPDF struct{}

func NewKnowYourCustomerPDF() *KnowYourClientPDF {
	return &KnowYourClientPDF{}
}

func (p KnowYourClientPDF) GenerateReport(data ConoceTuCliente) ([]byte, error) {

	tpl, err := utility.LoadTemplate("docs/conoceTuCliente.gohtml")
	if err != nil {
		return nil, err
	}

	html, err := utility.ExecuteTemplateToHTML(tpl, data)
	if err != nil {
		return nil, err
	}

	return utility.GeneratePDFFromHTML(html, utility.NewMargin(20, 20, 20, 20))
}

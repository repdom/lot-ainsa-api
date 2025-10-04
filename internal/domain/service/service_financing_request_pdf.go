package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/client/financing"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"time"

	"github.com/user0608/numeroaletras"
)

type FinancingRequestService struct {
	api             *financing.API
	pdf             *pdf.FinancingRequestPDF
	numberConverter *numeroaletras.NumeroALetras
}

func (s FinancingRequestService) GenerateReport(jwt, user, lang string, financingId int) ([]byte, *string, error) {
	load, err := s.api.LoadFinancing(jwt, user, lang, financingId)
	if err != nil {
		return nil, nil, err
	}

	fullName := load.Customer.Names + " " + load.Customer.LastNames

	por := validDataFloat(load.FinancingAmount) / (load.Lot.Price / 100)
	porT, _ := s.numberConverter.ToWords(por, 2)
	now := time.Now()

	data := pdf.SolicitudFinanciamiento{
		NombreCompleto:              fullName,
		CorreoElectronico:           validData(load.Customer.Email),
		Telefono:                    load.Customer.PhoneNumber,
		DireccionCompleta:           load.Customer.ResidentialAddress,
		Ciudad:                      validData(load.Customer.City),
		NumeroDUI:                   validData(load.Customer.Document.DUI),
		NombreUrbanizacion:          load.Lot.Development.Name,
		NombreDestinatario:          load.Lot.Development.OwnerName,
		MontoFinanciamiento:         fmt.Sprintf("%.2f", validDataFloat(load.FinancingAmount)),
		PorcentajeFinanciamiento:    fmt.Sprintf("%.2f", por),
		PorcentajeFinanciamientoTxt: porT,
		NumeroLote:                  load.Lot.Number,
		Poligono:                    load.Lot.Polygon,
		PrecioLote:                  fmt.Sprintf("%.2f", load.Lot.Price),
		AbonoEfectivo:               fmt.Sprintf("%.2f", validDataFloat(load.DownPaymentBalance)),
		Plazo:                       fmt.Sprintf("%d", load.TermElapsed),
		ActividadEconomica:          load.Customer.Financial.Occupation,
		Fecha:                       now.Format(dateFormat),
		Departamento:                "",
	}
	pdfData, fail := s.pdf.GenerateReport(data)
	return pdfData, &fullName, fail

}

func validDataFloat(data *float64) float64 {
	if data == nil {
		return 0.0
	}
	return *data
}

func NewFinancingRequestService(env *config.Env) port.FinancingRequestService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	return &FinancingRequestService{
		api:             api,
		pdf:             pdf.NewFinancingRequestPDF(),
		numberConverter: numeroaletras.NewNumeroALetras(),
	}
}

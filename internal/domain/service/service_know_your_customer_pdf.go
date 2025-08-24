package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/client/customer"
	"be-lotsanmateo-api/internal/adapter/externalapi/model"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/port"
	"log"
)

type ServiceKnowYourCustomer struct {
	api *customer.CustomerAPI
	pdf *pdf.KnowYourClientPDF
}

func (s ServiceKnowYourCustomer) GenerateReport(jwt, user, lang string, customerId int) ([]byte, *string, error) {
	log.Println("GenerateReport know your customer")
	customer, err := s.api.GetCustomerId(jwt, user, lang, customerId)
	if err != nil {
		return nil, nil, err
	}

	documentType, documentNumber := documents(customer)
	fullName := customer.Names + " " + customer.LastNames
	isPep := "No"
	if customer.Pep != nil && customer.Pep.Pep == true {
		isPep = "Si"
	}

	data := pdf.ConoceTuCliente{
		NombreCompleto:         fullName,
		TipoDocumento:          documentType,
		NumeroDocumento:        documentNumber,
		FechaNacimiento:        customer.Birthday,
		Genero:                 *customer.Gender.Detail,
		EstadoCivil:            customer.CivilStatus.CivilStatus,
		Nacionalidad:           *customer.Nationality,
		DireccionResidencia:    customer.ResidentialAddress,
		TelefonoContacto:       customer.PhoneNumber,
		CorreoElectronico:      *customer.Email,
		Ciudad:                 *customer.City,
		CodigoPostal:           *customer.ZipCode,
		OcupacionProfesion:     validData(customer.Profession),
		NombreEmpresa:          customer.Financial.EmployerName,
		CargoPuesto:            customer.Financial.Position,
		FuenteIngresos:         customer.Financial.IncomeSource,
		RangoIngresosMensuales: customer.Financial.EstimatedIncomeRange,
		PropositoRelacion:      customer.Financial.MainPurpose,
		EsPEP:                  isPep,
		DetallesPEP:            validData(customer.Pep.Details),
		Firma:                  fullName,
	}
	pdfData, fail := s.pdf.GenerateReport(data)
	return pdfData, &fullName, fail
}

func documents(customer *model.CustomerDomain) (documentType, documentNumber string) {
	documents := customer.Document
	if documents.DUI != nil {
		documentType = "DUI"
		documentNumber = *documents.DUI
	} else if documents.NIT != nil {
		documentType = "NIT"
		documentNumber = *documents.NIT
	} else if documents.Passport != nil {
		documentType = "PASAPORTE"
		documentNumber = *documents.Passport
	} else {
		documentType = "OTRO"
		documentNumber = ""
	}
	return documentType, documentNumber
}

func validData(data *string) string {
	if data == nil {
		return ""
	}
	return *data
}

func NewKnowYourCustomerService(env *config.Env) port.KnowYourCustomerService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := customer.NewCustomerClient(baseURL)
	return &ServiceKnowYourCustomer{
		api: api,
		pdf: pdf.NewKnowYourCustomerPDF(),
	}
}

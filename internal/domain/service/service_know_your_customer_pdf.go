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
	log.Println("GenerateReport know your customerDomain")
	customerDomain, err := s.api.GetCustomerId(jwt, user, lang, customerId)
	if err != nil {
		return nil, nil, err
	}

	documentType, documentNumber := documents(customerDomain)
	fullName := customerDomain.Names + " " + customerDomain.LastNames
	isPep := "No"
	if customerDomain.Pep != nil && customerDomain.Pep.Pep == true {
		isPep = "Si"
	}

	data := pdf.ConoceTuCliente{
		NombreCompleto:         fullName,
		TipoDocumento:          documentType,
		NumeroDocumento:        documentNumber,
		FechaNacimiento:        customerDomain.Birthday,
		Genero:                 *customerDomain.Gender.Detail,
		EstadoCivil:            customerDomain.CivilStatus.CivilStatus,
		Nacionalidad:           *customerDomain.Nationality,
		DireccionResidencia:    customerDomain.ResidentialAddress,
		TelefonoContacto:       customerDomain.PhoneNumber,
		CorreoElectronico:      *customerDomain.Email,
		Ciudad:                 *customerDomain.City,
		CodigoPostal:           *customerDomain.ZipCode,
		OcupacionProfesionDUI:  validData(customerDomain.Profession),
		OcupacionProfesion:     customerDomain.Financial.Occupation,
		NombreEmpresa:          customerDomain.Financial.EmployerName,
		CargoPuesto:            customerDomain.Financial.Position,
		FuenteIngresos:         customerDomain.Financial.IncomeSource,
		RangoIngresosMensuales: customerDomain.Financial.EstimatedIncomeRange,
		PropositoRelacion:      customerDomain.Financial.MainPurpose,
		EsPEP:                  isPep,
		DetallesPEP:            validData(customerDomain.Pep.Details),
		FullNamePep:            validData(customerDomain.Pep.FullName),
		TitlePep:               validData(customerDomain.Pep.Title),
		RelationshipPep:        validData(customerDomain.Pep.Relationship),
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

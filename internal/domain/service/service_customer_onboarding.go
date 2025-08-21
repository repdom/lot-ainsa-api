package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/client/customer"
	modelApi "be-lotsanmateo-api/internal/adapter/externalapi/model"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"
	"strings"
)

type customerOnboardingService struct {
	customerApi *customer.CustomerAPI
}

func (c customerOnboardingService) CreateCustomer(jwt, user, lang string, customerOnboarding model.RequestCustomerOnboarding) (error, error) {
	log.Println("Create domain")
	badRequest, internalServer := c.customerApi.ExistCustomer(jwt, user, lang, customerOnboarding.DocumentNumber, customerOnboarding.DocumentType)
	log.Println(internalServer, badRequest)
	if internalServer != nil || badRequest != nil {
		return badRequest, internalServer
	}

	document := modelApi.DocumentDomain{}

	switch strings.ToUpper(customerOnboarding.DocumentType) {
	case "DUI":
		document.DUI = &customerOnboarding.DocumentNumber
		break
	case "NIT":
		document.NIT = &customerOnboarding.DocumentNumber
		break
	case "PASAPORTE":
		document.Passport = &customerOnboarding.DocumentNumber
		break
	default:
		return nil, fmt.Errorf("el tipo de documento no es valido")
	}

	customerPep := modelApi.PepDomain{}
	if customerOnboarding.DetailsPep != "" {
		customerPep.Pep = true
		customerPep.Details = &customerOnboarding.DetailsPep
	}
	if customerOnboarding.FullName != "" {
		customerPep.Pep = true
		customerPep.FullName = &customerOnboarding.FullName
	}
	if customerOnboarding.Title != "" {
		customerPep.Pep = true
		customerPep.Title = &customerOnboarding.Title
	}
	if customerOnboarding.Relationship != "" {
		customerPep.Pep = true
		customerPep.Relationship = &customerOnboarding.Relationship
	}

	domain := modelApi.CustomerDomain{
		Names:              customerOnboarding.Names,
		LastNames:          customerOnboarding.Lastnames,
		Nationality:        customerOnboarding.Nationality,
		Document:           document,
		Gender:             modelApi.GenderDomain{Gender: customerOnboarding.Gender},
		CivilStatus:        modelApi.CivilStatusDomain{CivilStatus: customerOnboarding.MaritalStatus},
		ResidentialAddress: customerOnboarding.Address,
		Birthday:           customerOnboarding.BirthDate,
		City:               customerOnboarding.City,
		Email:              customerOnboarding.Email,
		PhoneNumber:        customerOnboarding.Phone,
		Pep:                customerPep,
		ZipCode:            customerOnboarding.PostalCode,
		Profession:         customerOnboarding.ProfessionDUI,
		Financial: modelApi.FinancialDomain{
			Position:             customerOnboarding.Position,
			EmployerName:         customerOnboarding.Employer,
			EstimatedIncomeRange: customerOnboarding.RangeIncome,
			IncomeSource:         customerOnboarding.SourceOfIncome,
			Occupation:           customerOnboarding.Occupation,
			MainPurpose:          customerOnboarding.RelationFinancial,
		},
	}

	badRequest, internalServer = c.customerApi.CreateCustomer(jwt, user, lang, domain)
	log.Println(internalServer, "|", badRequest)
	if internalServer != nil || badRequest != nil {
		return badRequest, internalServer
	}

	return nil, nil
}

func NewCustomerOnboardingService(env *config.Env) port.CustomerOnboardingService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "http://localhost:8080")
	customerApi := customer.NewCustomerClient(baseURL)
	return &customerOnboardingService{
		customerApi: customerApi,
	}
}

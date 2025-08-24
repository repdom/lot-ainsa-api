package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/client/financing"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/port"
)

type FinancingRequestService struct {
	api *financing.API
	pdf *pdf.FinancingRequestPDF
}

func (s FinancingRequestService) GenerateReport(jwt, user, lang string, financingId int) ([]byte, *string, error) {
	//TODO implement me
	panic("implement me")
}

func NewFinancingRequestService(env *config.Env) port.FinancingRequestService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	return &FinancingRequestService{
		api: api,
		pdf: pdf.NewFinancingRequestPDF(),
	}
}

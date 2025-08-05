package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/financing"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/port"
)

type ServiceInvoice struct {
	api *financing.API
}

func (s ServiceInvoice) InvoicePayment(paymentId int) ([]byte, error) {
	panic("implement me")
}

func (s ServiceInvoice) InvoiceReservation(reservationId int) ([]byte, error) {
	panic("implement me")
}

func (s ServiceInvoice) InvoicePremium(premiumId int) ([]byte, error) {
	panic("implement me")
}

func NewInvoiceService(env *config.Env) port.InvoiceService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	return &ServiceInvoice{
		api: api,
	}
}

package pkg

import (
	"be-lotsanmateo-api/internal/adapter/http/handler/api"
	"be-lotsanmateo-api/internal/adapter/http/handler/pdf"
	"be-lotsanmateo-api/internal/config"
	service "be-lotsanmateo-api/internal/domain/service"
)

func NewCustomerHandler(env *config.Env) *api.CustomerOnboardingHandler {
	servicePort := service.NewCustomerOnboardingService(env)
	return api.NewCustomerOnboardingHandler(servicePort)
}

func NewPaymentPlanHandler() *api.PaymentPlanSimulationHandler {
	servicePort := service.NewCalculatePlan()
	return api.NewPaymentSimulationHandler(servicePort)
}

func NewCalculatePlanPdfHandler(env *config.Env) *pdf.PaymentPlanDocumentHandler {
	servicePort := service.NewCalculatePlanPDF(env)
	return pdf.NewPaymentPlanDocumentHandler(servicePort)
}

func NewPagareHandler(env *config.Env) *pdf.PagareHandler {
	pagare := service.NewPagarePDF(env)
	return pdf.NewPagareHandler(pagare)
}

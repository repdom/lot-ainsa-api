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

func NewPromissoryNoteHandler(env *config.Env) *pdf.PromissoryNoteHandler {
	promissoryNote := service.NewPromissoryNotePDF(env)
	return pdf.NewPromissoryNoteHandler(promissoryNote)
}

func NewLoanPaymentHandler(env *config.Env) *api.LoanPaymentHandler {
	loanPayment := service.NewLoanPaymentService(env)
	return api.NewLoanPaymentHandler(loanPayment)
}

func NewFinancingHandler(env *config.Env) *api.FinancingsActiveHandler {
	activated := service.NewServiceFinancingsActions(env)
	return api.NewFinancingHandler(activated)
}

func NewInvoiceHandler(env *config.Env) *pdf.InvoiceDocumentHandler {
	invoice := service.NewInvoiceService(env)
	return pdf.NewInvoiceDocumentHandler(invoice)
}

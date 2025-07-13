package pkg

import (
	"be-lotsanmateo-api/internal/adapter/http/handler/api"
	"be-lotsanmateo-api/internal/config"
	service "be-lotsanmateo-api/internal/domain/service"
)

func NewCustomerHandler(env *config.Env) *api.CustomerOnboardingHandler {
	servicePort := service.NewCustomerOnboardingService(env)
	return api.NewCustomerOnboardingHandler(servicePort)
}

func NewPaymentPlanHandler(env *config.Env) *api.PaymentPlanSimulationHandler {
	servicePort := service.NewCalculatePlan(env)
	return api.NewPaymentSimulationHandler(servicePort)
}

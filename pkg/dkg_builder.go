package pkg

import (
	"be-lotsanmateo-api/internal/adapter/http/handler/api"
	"be-lotsanmateo-api/internal/config"
	service "be-lotsanmateo-api/internal/domain/service"
)

func NewCustomerHandler(env *config.Env) *api.CustomerOnboardingHandler {
	servicePort := service.NewCustomerOnboardingService(env)
	customer := api.NewCustomerOnboardingHandler(servicePort)
	return customer
}

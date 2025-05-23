package port

import "be-lotsanmateo-api/internal/domain/model"

type CustomerOnboardingService interface {
	CreateCustomer(jwt, user, lang string, customer model.RequestCustomerOnboarding) (error, error)
}

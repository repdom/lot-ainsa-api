package port

import "be-lotsanmateo-api/internal/domain/model"

type FinancingActionService interface {
	Activation(jwt, user, lang string, loan model.RequestLoan, financingId int) error
}

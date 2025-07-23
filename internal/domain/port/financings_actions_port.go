package port

import "be-lotsanmateo-api/internal/domain/model"

type FinancingsActionService interface {
	Activation(loan model.RequestLoan, financingId int) error
}

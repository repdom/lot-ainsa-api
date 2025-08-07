package port

import "be-lotsanmateo-api/internal/domain/model"

type LoanPaymentService interface {
	CalculateLoanPayment(jwt, user, lang string, financingId int, share float64) (*model.PaymentLoanResponse, error)
}

type ApiService interface {
	GenerateSimulation(request model.RequestLoan) (*model.ResponseLoan, error)
}

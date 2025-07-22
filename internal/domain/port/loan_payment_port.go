package port

import "be-lotsanmateo-api/internal/domain/model"

type LoanPaymentService interface {
	CalculateLoanPayment(financingId int, share float64) (*model.PaymentLoanResponse, error)
}

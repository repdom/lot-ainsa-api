package port

import "be-lotsanmateo-api/internal/domain/model"

type ReportService interface {
	GenerateReport(financingId int) ([]byte, error)
}

type ApiService interface {
	GenerateSimulation(request model.RequestLoan) (*model.ResponseLoan, error)
}

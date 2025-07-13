package port

import "be-lotsanmateo-api/internal/domain/model"

type ReportSimulationService interface {
	GenerateReport(idLot int) ([]byte, error)
}

type ReportService interface {
	GenerateReport(idLot, customerId int) ([]byte, error)
}

type ApiService interface {
	Generate(idFinancial string) (*model.ResponseLoan, error)

	GenerateSimulation(request model.RequestLoan) (*model.ResponseLoan, error)
}

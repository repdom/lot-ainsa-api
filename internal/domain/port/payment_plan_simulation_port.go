package port

import "be-lotsanmateo-api/internal/domain/model"

type PdfGenerator interface {
	Generate(data model.ResponseLoan) ([]byte, error)
}

type ReportSimulationService interface {
	GenerateReport(idLot int) ([]byte, error)
}

type ReportService interface {
	GenerateReport(idLot, customerId int) ([]byte, error)
}

type ApiService interface {
	GenerateReport(idLot int) (model.ResponseLoan, error)
}

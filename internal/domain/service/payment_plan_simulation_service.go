package service

import "be-lotsanmateo-api/internal/domain/port"

type paymentPlanSimulationService struct {
	pdfGen port.PdfGenerator
}

func (p paymentPlanSimulationService) GenerateReport(idLot int) ([]byte, error) {
	panic("implement me")
}

func NewPaymentPlanSimulationService(pdfGen port.PdfGenerator) port.ReportSimulationService {
	return &paymentPlanSimulationService{pdfGen}
}

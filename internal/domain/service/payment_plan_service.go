package service

import "be-lotsanmateo-api/internal/domain/port"

type paymentPlanService struct {
	pdfGen port.PdfGenerator
}

func (p paymentPlanService) GenerateReport(idLot, customerId int) ([]byte, error) {
	panic("implement me")
}

func NewPaymentPlanService(pdfGen port.PdfGenerator) port.ReportService {
	return &paymentPlanService{pdfGen}
}

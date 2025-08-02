package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/financing"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"log"
	"time"
)

type PaymentPlanPDF struct {
	api           *financing.API
	calculatePlan *CalculatePlan
}

func (p PaymentPlanPDF) GenerateReport(financingId int) ([]byte, error) {
	loadFinancing, err := p.api.LoadFinancing("", "", "", financingId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	request := model.RequestLoan{}
	request.Rate = loadFinancing.InterestRate
	request.Amount = loadFinancing.Amount
	log.Printf("Total term: %d", loadFinancing.TotalTerm)
	if loadFinancing.TotalTerm == 0 {
		request.Months = 240
	} else {
		request.Months = loadFinancing.TotalTerm
	}
	request.Premium = loadFinancing.DownPaymentRate
	if loadFinancing.StartDate != "" {
		t, _ := time.Parse("2006-01-02", loadFinancing.StartDate)
		request.Payday = t.Day()
	} else {
		request.Payday = time.Now().Day()
	}

	land, err := p.calculatePlan.GenerateSimulation(request)
	if err != nil {
		return nil, err
	}

	paymentPDF := pdf.NewPaymentPlanPDF()

	var paymentPlan pdf.PaymentPlan
	paymentPlan.Loan = *land
	if loadFinancing.StartDate != "" {
		t, _ := time.Parse("2006-01-02", loadFinancing.StartDate)
		paymentPlan.StartDate = t.Format("02/01/2006")
	} else {
		paymentPlan.StartDate = time.Now().Format("02/01/2006")
	}
	paymentPlan.DUI = loadFinancing.Customer.Document.DUI
	paymentPlan.FullName = loadFinancing.Customer.Names + " " + loadFinancing.Customer.LastNames
	paymentPlan.Address = loadFinancing.Customer.ResidentialAddress
	paymentPlan.Phone = loadFinancing.Customer.PhoneNumber
	paymentPlan.Lote = loadFinancing.Lot.Number
	paymentPlan.Polygon = loadFinancing.Lot.Polygon
	paymentPlan.Area = loadFinancing.Lot.Area

	return paymentPDF.GeneratePDF(paymentPlan)
}

func NewCalculatePlanPDF(env *config.Env) port.ReportService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	return &PaymentPlanPDF{
		api:           api,
		calculatePlan: &CalculatePlan{},
	}
}

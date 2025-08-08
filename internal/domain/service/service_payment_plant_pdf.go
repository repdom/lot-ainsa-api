package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/client/financing"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"
	"time"
)

type PaymentPlanPDF struct {
	api           *financing.API
	calculatePlan *CalculatePlan
}

func (p PaymentPlanPDF) GenerateReport(jwt, user, lang string, financingId int) ([]byte, *string, error) {
	loadFinancing, err := p.api.LoadFinancing(jwt, user, lang, financingId)
	if err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}
	request := model.RequestLoan{}
	if loadFinancing.InterestRate == nil {
		log.Print("el financiamiento no tiene interes")
		return nil, nil, fmt.Errorf("el financiamiento no tiene interes")
	}
	request.Rate = *loadFinancing.InterestRate
	request.Amount = loadFinancing.Amount
	log.Printf("Total term: %d", loadFinancing.TotalTerm)
	if loadFinancing.TotalTerm == nil {
		request.Months = 240
	}
	request.Months = *loadFinancing.TotalTerm
	if loadFinancing.DownPaymentBalance != nil {
		request.Premium = *loadFinancing.DownPaymentBalance
	} else {
		request.Premium = 0.0
	}
	if loadFinancing.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *loadFinancing.StartDate)
		request.Payday = t.Day()
	} else {
		request.Payday = time.Now().Day()
	}

	land, err := p.calculatePlan.GenerateSimulation(request)
	if err != nil {
		return nil, nil, err
	}

	paymentPDF := pdf.NewPaymentPlanPDF()

	var paymentPlan pdf.PaymentPlan
	paymentPlan.Loan = *land
	if loadFinancing.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *loadFinancing.StartDate)
		paymentPlan.StartDate = t.Format("02/01/2006")
	} else {
		paymentPlan.StartDate = time.Now().Format("02/01/2006")
	}

	if loadFinancing.Customer.Document.DUI != nil {
		paymentPlan.DUI = *loadFinancing.Customer.Document.DUI
	}

	if loadFinancing.Customer.Document.NIT != nil {
		paymentPlan.DUI = *loadFinancing.Customer.Document.NIT
	}

	if loadFinancing.Customer.Document.Passport != nil {
		paymentPlan.DUI = *loadFinancing.Customer.Document.Passport
	}

	paymentPlan.FullName = loadFinancing.Customer.Names + " " + loadFinancing.Customer.LastNames
	paymentPlan.Address = loadFinancing.Customer.ResidentialAddress
	paymentPlan.Phone = loadFinancing.Customer.PhoneNumber
	paymentPlan.Lote = loadFinancing.Lot.Number
	paymentPlan.Polygon = loadFinancing.Lot.Polygon
	paymentPlan.Area = loadFinancing.Lot.Area

	clientName := loadFinancing.Customer.Names + " " + loadFinancing.Customer.LastNames

	var pdfData []byte
	pdfData, fail := paymentPDF.GeneratePDF(paymentPlan)
	return pdfData, &clientName, fail
}

func NewCalculatePlanPDF(env *config.Env) port.ReportService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	return &PaymentPlanPDF{
		api:           api,
		calculatePlan: &CalculatePlan{},
	}
}

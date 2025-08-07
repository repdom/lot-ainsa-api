package service

import (
	financingApi "be-lotsanmateo-api/internal/adapter/externalapi/client/financing"
	model2 "be-lotsanmateo-api/internal/adapter/externalapi/model"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"log"
	"time"
)

type ServiceFinancingsActions struct {
	calculatePlan port.ApiService
	api           *financingApi.API
}

func (s ServiceFinancingsActions) Activation(loan model.RequestLoan, financingId int) error {
	cal, err := s.calculatePlan.GenerateSimulation(loan)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	log.Printf("%f", cal.Amount)
	var financings model2.FinancingDomain

	financings.Amount = cal.Amount
	financings.FinancingAmount = &cal.TotalAmount
	financings.FinancingBalance = 0.0
	financings.FinancingAmountPending = &cal.TotalAmount
	financings.DownPaymentRate = &cal.DownPaymentRate
	startDate := time.Now().Format("2006-01-02")
	log.Print("fecha de inicio: ", startDate)
	financings.StartDate = &startDate
	financings.TotalTerm = &cal.NumberOfInstallments
	financings.TermElapsed = 0
	financings.MissingTerm = &cal.NumberOfInstallments
	financings.MonthlyPayment = &cal.MonthlyPayment
	financings.InterestRate = &cal.Rate
	financings.InterestRateMonthly = &cal.RateMonths
	financings.Status = "active"
	patchFinancing, err := s.api.PatchFinancing("", "", "", financingId, financings)
	if err != nil {
		return err
	}
	if patchFinancing != nil {
		log.Printf("financings id pathc: %d", patchFinancing.ID)
	} else {
		log.Printf("financings id pathc: %d", financingId)
	}
	return nil
}

func NewServiceFinancingsActions(env *config.Env) port.FinancingsActionService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financingApi.NewFinancingAPI(baseURL)
	return &ServiceFinancingsActions{
		calculatePlan: NewCalculatePlan(),
		api:           api,
	}
}

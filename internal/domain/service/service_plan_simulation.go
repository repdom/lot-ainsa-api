package service

import (
	"be-lotsanmateo-api/internal/domain"
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"math"
	"time"
)

type CalculatePlan struct {
}

func (c CalculatePlan) GenerateSimulation(request model.RequestLoan) (*model.ResponseLoan, error) {
	response := model.ResponseLoan{}
	response.Rate = request.Rate

	response.Years = float64(request.Months / 12)
	response.Amount = request.Amount

	n := float64(request.Months)
	p := request.Amount

	porce := request.Amount / 100
	premiun := request.Premium / porce

	response.PremiumRate = premiun
	response.DownPaymentRate = premiun

	p = p - request.Premium

	response.TotalAmount = p
	response.Premium = request.Premium

	if request.Rate == 0 {
		payment := p / n
		response.MonthlyPayment = domain.RoundToTwoDecimals(payment)
		rateMont := 0.0
		response.RateMonths = domain.RoundToTwoDecimals(rateMont)
		response.NumberOfInstallments = request.Months
		response.InterestsTotal = 0
		response.TotalPayments = domain.RoundToTwoDecimals(request.Amount + request.Premium)
		response.FeeSimulation = FeeSimulationCalculate(request.Months, payment, 0.0, p, request)
		return &response, nil
	}

	rateMonths := (request.Rate / 100) / 12

	rateMont := request.Rate / 12
	response.RateMonths = domain.RoundToTwoDecimals(rateMont)
	response.NumberOfInstallments = request.Months

	payment := (p * rateMonths) / (1 - math.Pow(1+rateMonths, -n))
	paymentR := domain.RoundToTwoDecimals(payment)
	response.MonthlyPayment = paymentR
	interestTotal := domain.RoundToTwoDecimals((payment * n) - p)
	response.InterestsTotal = interestTotal
	response.TotalPayments = domain.RoundToTwoDecimals(interestTotal + p + request.Premium)
	response.FeeSimulation = FeeSimulationCalculate(request.Months, payment, rateMonths, p, request)

	return &response, nil
}

func FeeSimulationCalculate(Months int, payment float64, rateMonths, p float64, loan model.RequestLoan) []model.FreeSimulation {
	var feeSimulations []model.FreeSimulation
	balance := p
	start := time.Now()
	dayOfMonth := loan.Payday

	for month := 1; month <= Months; month++ {
		year, m := start.AddDate(0, month, 0).Year(), start.AddDate(0, month, 0).Month()
		firstOfNextMonth := time.Date(year, m+1, 1, 0, 0, 0, 0, time.UTC)
		lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
		day := dayOfMonth
		if dayOfMonth > lastOfMonth.Day() {
			day = lastOfMonth.Day()
		}
		paymentDate := time.Date(year, m, day, 0, 0, 0, 0, time.UTC)

		var feeSimulation model.FreeSimulation
		feeSimulation.Payday = paymentDate.Format("02/01/2006")
		feeSimulation.BalanceStart = domain.RoundToTwoDecimals(balance)
		interest := balance * rateMonths
		principal := payment - interest
		balance = balance - principal
		feeSimulation.Amount = domain.RoundToTwoDecimals(payment)
		feeSimulation.BalanceLast = domain.RoundToTwoDecimals(balance)
		feeSimulation.Capital = domain.RoundToTwoDecimals(principal)
		feeSimulation.Interest = domain.RoundToTwoDecimals(interest)

		feeSimulations = append(feeSimulations, feeSimulation)
	}

	return feeSimulations
}

func NewCalculatePlan() port.ApiService {
	return &CalculatePlan{}
}

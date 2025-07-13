package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/financing"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"log"
	"math"
	"time"
)

type CalculatePlan struct {
	api *financing.API
}

func (c CalculatePlan) Generate(idFinancial string) (*model.ResponseLoan, error) {
	loadFinancing, err := c.api.LoadFinancing("", "", "", idFinancial)
	if err != nil {
		return nil, err
	}
	request := model.RequestLoan{}
	request.Rate = loadFinancing.InterestRate
	request.Amount = loadFinancing.Amount
	request.Months = loadFinancing.TotalTerm
	request.Premium = loadFinancing.Amount - loadFinancing.FinancingAmount
	request.Payday = loadFinancing.StartDate.Day()
	return c.GenerateSimulation(request)

}

func (c CalculatePlan) GenerateSimulation(request model.RequestLoan) (*model.ResponseLoan, error) {
	response := model.ResponseLoan{}

	n := float64(request.Months)
	p := request.Amount

	response.PremiumRate = request.Premium
	premium := request.Premium / 100

	pm := p * premium

	p = p - pm

	response.TotalAmount = p
	response.Premium = pm

	if request.Rate == 0 {
		payment := p / n
		response.MonthlyPayment = roundToTwoDecimals(payment)
		rateMont := 0.0
		response.RateMonths = roundToTwoDecimals(rateMont)
		response.NumberOfInstallments = request.Months
		response.InterestsTotal = 0
		response.TotalPayments = roundToTwoDecimals(request.Amount + pm)
		response.FeeSimulation = FeeSimulationCalculate(request.Months, payment, 0.0, p, request)
		return &response, nil
	}

	rateMonths := (request.Rate / 100) / 12

	rateMont := request.Rate / 12
	response.RateMonths = roundToTwoDecimals(rateMont)
	response.NumberOfInstallments = request.Months

	payment := (p * rateMonths) / (1 - math.Pow(1+rateMonths, -n))
	paymentR := roundToTwoDecimals(payment)
	response.MonthlyPayment = paymentR
	log.Printf("paymentR = %f\n", (payment*n)-p)
	interestTotal := roundToTwoDecimals((payment * n) - p)
	response.InterestsTotal = interestTotal
	response.TotalPayments = roundToTwoDecimals(interestTotal + p + pm)
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
		feeSimulation.BalanceStart = roundToTwoDecimals(balance)
		interest := balance * rateMonths
		principal := payment - interest
		balance = balance - principal
		feeSimulation.Amount = roundToTwoDecimals(payment)
		feeSimulation.BalanceLast = roundToTwoDecimals(balance)
		feeSimulation.Capital = roundToTwoDecimals(principal)
		feeSimulation.Interest = roundToTwoDecimals(interest)

		feeSimulations = append(feeSimulations, feeSimulation)
	}

	return feeSimulations
}

func roundToTwoDecimals(value float64) float64 {
	var round float64
	if ((value * 100) - math.Floor(value*100)) >= 0.5 {
		round = math.Ceil(value*100) / 100
	} else {
		round = math.Floor(value*100) / 100
	}
	log.Printf("roundToTwoDecimals(%f) \n", value)
	log.Println("round = ", round)
	return round
}

func NewCalculatePlan(env *config.Env) port.ApiService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "http://localhost:8080")
	api := financing.NewFinancingAPI(baseURL)
	return &CalculatePlan{
		api: api,
	}
}

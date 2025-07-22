package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/financing"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain"
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"
	"sort"
	"time"
)

type ServiceLoanPayment struct {
	api *financing.API
}

func (s ServiceLoanPayment) CalculateLoanPayment(financingId int, share float64) (*model.PaymentLoanResponse, error) {
	loadFinancing, err := s.api.LoadFinancing("", "", "", financingId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Printf(loadFinancing.StartDate)

	amount := loadFinancing.Balance
	rate := loadFinancing.InterestRateMonthly / 100
	now := time.Now()

	var paymentLastDate time.Time
	if loadFinancing.Payments != nil && len(loadFinancing.Payments) > 0 {
		payments := loadFinancing.Payments
		slice := payments[:]
		sort.Slice(slice, func(i, j int) bool {
			dateI := slice[i].PaymentDate
			dateJ := slice[j].PaymentDate
			return dateJ.Before(dateI)
		})
		paymentLastDate = slice[0].PaymentDate
	} else {
		paymentLastDate, err = time.Parse("2006-01-02", loadFinancing.StartDate)
		if err != nil {
			log.Printf(err.Error())
			return nil, fmt.Errorf("no se pudo formatear la fecha: '%s'", loadFinancing.StartDate)
		}
	}

	days := int(now.Sub(paymentLastDate).Hours() / 24)

	var penalty float64
	if days > 30 {
		penalty = 15.00
	} else {
		penalty = 0.0
	}

	interest := amount * rate * float64(days) / 30.0

	var capital float64

	if share != 0 {
		capital = share - interest
	} else {
		capital = loadFinancing.MonthlyPayment - interest
	}

	var shareAmount = capital + interest + penalty

	var balanceStart = loadFinancing.Balance
	var lastBalance = balanceStart - shareAmount

	return &model.PaymentLoanResponse{
		Interest:      domain.RoundToTwoDecimals(interest),
		Share:         domain.RoundToTwoDecimals(shareAmount),
		Capital:       domain.RoundToTwoDecimals(capital),
		Penalty:       domain.RoundToTwoDecimals(penalty),
		AmountBalance: domain.RoundToTwoDecimals(lastBalance),
		AmountStart:   domain.RoundToTwoDecimals(balanceStart),
		Customer:      loadFinancing.Customer,
		Lot:           loadFinancing.Lot,
	}, nil

}

func NewLoanPaymentService(env *config.Env) port.LoanPaymentService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	return &ServiceLoanPayment{
		api: api,
	}
}

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

	amount := loadFinancing.FinancingAmount
	rate := loadFinancing.InterestRateMonthly / 100

	log.Print("amount: ", amount)
	log.Print("rate: ", rate)
	log.Print("rate year:", loadFinancing.InterestRate)

	now := time.Now()

	var paymentLastDate time.Time
	if loadFinancing.Payments != nil && len(loadFinancing.Payments) > 0 {
		log.Print("el financiamiento cuenta con pagos ")
		payments := loadFinancing.Payments
		slice := payments[:]
		sort.Slice(slice, func(i, j int) bool {
			dateI := slice[i].PaymentDate
			dateJ := slice[j].PaymentDate
			return dateJ.Before(dateI)
		})
		paymentLastDate = slice[0].PaymentDate
	} else {
		log.Print("el financiamiento no cuenta con pagos")
		paymentLastDate, err = time.Parse("2006-01-02", loadFinancing.StartDate)
		if err != nil {
			log.Printf(err.Error())
			return nil, fmt.Errorf("no se pudo formatear la fecha: '%s'", loadFinancing.StartDate)
		}
	}

	log.Print("Ultima fetch de pago: ", paymentLastDate)

	days := int(now.Sub(paymentLastDate).Hours() / 24)

	log.Print(days)
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

	var lastBalance = amount - ((shareAmount - penalty) + interest)

	if lastBalance < 0 {
		capital = capital + lastBalance
		shareAmount = capital + interest + penalty
		lastBalance = amount - ((shareAmount - penalty) + interest)
	}

	return &model.PaymentLoanResponse{
		Interest:      domain.RoundToTwoDecimals(interest),
		Share:         domain.RoundToTwoDecimals(shareAmount),
		Capital:       domain.RoundToTwoDecimals(capital),
		Penalty:       domain.RoundToTwoDecimals(penalty),
		AmountBalance: domain.RoundToTwoDecimals(lastBalance),
		AmountStart:   domain.RoundToTwoDecimals(amount),
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

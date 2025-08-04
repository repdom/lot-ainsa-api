package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/financing"
	modelFinancing "be-lotsanmateo-api/internal/adapter/externalapi/model"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain"
	"be-lotsanmateo-api/internal/domain/model"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"
	"sort"
	"time"
)

const (
	dateLayout       = "2006-01-02"
	timeLayout       = "2006-01-02T15:04:05"
	daysPerMonth     = 30.0
	daysPerYear      = 365
	penaltyThreshold = 30
)

type ServiceLoanPayment struct {
	api *financing.API
}

type validationError struct {
	field string
	msg   string
}

func (e validationError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.msg)
}

func (s ServiceLoanPayment) CalculateLoanPayment(financingId int, share float64) (*model.PaymentLoanResponse, error) {
	loadFinancing, err := s.api.LoadFinancing("", "", "", financingId)
	if err != nil {
		return nil, fmt.Errorf("error loading financing: %w", err)
	}

	// Validar datos requeridos
	if err := s.validateFinancingData(loadFinancing); err != nil {
		return nil, err
	}

	// Extraer valores validados
	amount := *loadFinancing.FinancingAmountPending
	rate := *loadFinancing.InterestRateMonthly / 100

	log.Printf("Processing financing %d - Amount: %.2f, Rate: %.4f", financingId, amount, rate)

	// Obtener fecha del último pago
	paymentLastDate, err := s.getPaymentDate(loadFinancing)
	if err != nil {
		return nil, err
	}

	// Calcular días transcurridos
	days := s.calculateDays(time.Now(), paymentLastDate)
	log.Printf("Days since last payment: %d", days)

	// Calcular componentes del pago
	interestOnArrears := loadFinancing.Lot.Development.InterestOnArrears
	interestOnArrearsAmount := s.calculatePenalty(days, amount, interestOnArrears)
	interest := s.calculateInterest(amount, rate, days)
	capital, err := s.calculateCapital(share, interest, loadFinancing.MonthlyPayment, interestOnArrears)

	if err != nil {
		return nil, err
	}

	// Calcular montos finales
	shareAmount := capital + interest + interestOnArrearsAmount
	lastBalance := amount - capital

	// Ajustar si el balance es negativo
	if lastBalance < 0 {
		capital += lastBalance
		shareAmount = capital + interest + interestOnArrearsAmount
		lastBalance := amount - capital
		if lastBalance < 0 {
			log.Print("Warning: balance is negative, using interest as capital")
			lastBalance = 0
		}
	}

	return &model.PaymentLoanResponse{
		Interest:      domain.RoundToTwoDecimals(interest),
		Share:         domain.RoundToTwoDecimals(shareAmount),
		Capital:       domain.RoundToTwoDecimals(capital),
		Penalty:       domain.RoundToTwoDecimals(interestOnArrearsAmount),
		AmountBalance: domain.RoundToTwoDecimals(lastBalance),
		AmountStart:   domain.RoundToTwoDecimals(amount),
		Customer:      loadFinancing.Customer,
		Lot:           loadFinancing.Lot,
	}, nil
}

func (s ServiceLoanPayment) validateFinancingData(financing *modelFinancing.FinancingDomain) error {
	if financing.StartDate == nil {
		return validationError{"StartDate", "el financiamiento no tiene fecha de inicio"}
	}
	if financing.FinancingAmountPending == nil {
		return validationError{"FinancingAmountPending", "el financiamiento no tiene monto pendiente"}
	}
	if financing.InterestRateMonthly == nil {
		return validationError{"InterestRateMonthly", "el financiamiento no tiene interés mensual"}
	}

	return nil
}

func (s ServiceLoanPayment) getPaymentDate(financing *modelFinancing.FinancingDomain) (time.Time, error) {
	if financing.Payments != nil && len(*financing.Payments) > 0 {
		log.Print("Financing has payments")
		return getLastPaymentDate(*financing.Payments)
	}

	log.Print("Financing has no payments, using start date")
	paymentDate, err := time.Parse(dateLayout, *financing.StartDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid start date format '%s': %w", *financing.StartDate, err)
	}
	return paymentDate, nil
}

func (s ServiceLoanPayment) calculateDays(now, lastPayment time.Time) int {
	return int(now.Sub(lastPayment).Hours() / 24)
}

func (s ServiceLoanPayment) calculatePenalty(days int, amount float64, interestOnArrears float64) float64 {
	if days > penaltyThreshold && interestOnArrears > 0 {
		interestOnArrearsDay := float64(days) / daysPerYear
		return amount * interestOnArrearsDay * float64(days) // No penalty applied at the moment
	}
	return 0.0
}

func (s ServiceLoanPayment) calculateInterest(amount, rate float64, days int) float64 {
	return amount * rate * float64(days) / daysPerMonth
}

func (s ServiceLoanPayment) calculateCapital(share, interest float64, monthlyPayment *float64, interestOnArrears float64) (float64, error) {
	if share != 0 {
		principal := share - (interest + interestOnArrears)
		if principal < 0 {
			log.Print("Warning: principal is negative, using interest as capital")
			return 0.0, fmt.Errorf("no se puede recibir un abono menor de los intereses")
		}
		return principal, nil
	}

	if monthlyPayment == nil {
		log.Print("Warning: monthly payment is nil, using interest as capital")
		return interest, fmt.Errorf("no se cuenta una cuota configurada")
	}

	principal := *monthlyPayment - interest
	if principal < 0 {
		log.Print("Warning: principal is negative, using interest as capital")
		return 0.0, fmt.Errorf("no se puede recibir un pago menor de los intereses")
	}

	return principal, nil
}

func getLastPaymentDate(payments []modelFinancing.PaymentDomain) (time.Time, error) {
	if len(payments) == 0 {
		return time.Time{}, fmt.Errorf("no payments available")
	}

	// Crear una copia para evitar modificar el slice original
	sortedPayments := make([]modelFinancing.PaymentDomain, len(payments))
	copy(sortedPayments, payments)

	// Ordenar por fecha descendente
	sort.Slice(sortedPayments, func(i, j int) bool {
		if sortedPayments[i].PaymentDate == nil || sortedPayments[j].PaymentDate == nil {
			return sortedPayments[i].PaymentDate != nil
		}

		dateI, errI := time.Parse(timeLayout, *sortedPayments[i].PaymentDate)
		dateJ, errJ := time.Parse(timeLayout, *sortedPayments[j].PaymentDate)

		if errI != nil || errJ != nil {
			return errI == nil
		}

		return dateI.After(dateJ)
	})

	if sortedPayments[0].PaymentDate == nil {
		return time.Time{}, fmt.Errorf("latest payment has no date")
	}

	return time.Parse(timeLayout, *sortedPayments[0].PaymentDate)
}

func NewLoanPaymentService(env *config.Env) port.LoanPaymentService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	return &ServiceLoanPayment{
		api: api,
	}
}

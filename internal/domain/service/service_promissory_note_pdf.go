package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/client/financing"
	"be-lotsanmateo-api/internal/adapter/externalapi/model"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/user0608/numeroaletras"
)

const (
	dateFormat        = "2006-01-02"
	displayDateFormat = "02/01/2006"
	currency          = "Dólares"
)

type PromissoryNotePDF struct {
	api               *financing.API
	generatePagarePDF pdf.GeneratePagarePDF
	numberConverter   *numeroaletras.NumeroALetras
}

func (p *PromissoryNotePDF) GenerateReport(jwt, user, lang string, financingId int) ([]byte, *string, error) {
	loadFinancing, err := p.api.LoadFinancing(jwt, user, lang, financingId)
	if err != nil {
		log.Printf("Error loading financing %d: %v", financingId, err)
		return nil, nil, fmt.Errorf("error cargando financiamiento: %w", err)
	}

	clientName := loadFinancing.Customer.Names + " " + loadFinancing.Customer.LastNames

	data, err := p.buildPagareData(loadFinancing)
	if err != nil {
		return nil, nil, err
	}

	var pdfData []byte
	pdfData, fail := p.generatePagarePDF.GenerateReport(data)
	return pdfData, &clientName, fail
}

func (p *PromissoryNotePDF) buildPagareData(financing *model.FinancingDomain) (pdf.PagareData, error) {
	var data pdf.PagareData

	// Validar y calcular edad del cliente
	age, err := p.calculateAge(financing.Customer.Birthday)
	if err != nil {
		return data, fmt.Errorf("error calculando edad del cliente: %w", err)
	}

	// Obtener documento de identidad
	identityDoc, err := p.getIdentityDocument(financing.Customer.Document)
	if err != nil {
		return data, err
	}

	// Validar campos requeridos
	if err := p.validateRequiredFields(financing); err != nil {
		return data, err
	}

	// Construir estructura de datos
	data = pdf.PagareData{
		NoteCode:              fmt.Sprintf("%09d", financing.ID),
		ClientName:            fmt.Sprintf("%s %s", financing.Customer.Names, financing.Customer.LastNames),
		IdentityDocument:      identityDoc,
		ClientAge:             age,
		Area:                  fmt.Sprintf("%.0f", financing.Lot.Area),
		LotNumber:             financing.Lot.Number,
		Address:               financing.Customer.ResidentialAddress,
		Block:                 financing.Lot.Polygon,
		TermInYears:           *financing.TotalTerm / 12,
		InstallmentCount:      *financing.TotalTerm,
		DateNow:               time.Now().Format(displayDateFormat),
		InterestRate:          fmt.Sprintf("%.2f", *financing.InterestRate),
		AmountToFinance:       formatAmount(*financing.FinancingAmount),
		InstallmentAmount:     formatAmount(*financing.MonthlyPayment),
		DownPaymentPercentage: fmt.Sprintf("%.2f", *financing.DownPaymentRate),
		InterestOnArrears:     fmt.Sprintf("%.2f", financing.Lot.Development.InterestOnArrears),
		LotCost:               formatAmount(financing.Amount),
		Profession:            financing.Customer.Financial.Occupation,
	}

	// Manejar pago inicial
	if financing.DownPaymentBalance != nil {
		data.DownPayment = formatAmount(*financing.DownPaymentBalance)
	} else {
		data.DownPaymentPercentage = "0.00"
		data.DownPayment = "0.00"
	}

	// Calcular fecha del primer pago
	if financing.StartDate != nil {
		startDate, err := time.Parse(dateFormat, *financing.StartDate)
		if err != nil {
			return data, fmt.Errorf("error parseando fecha de inicio: %w", err)
		}
		data.FirstPaymentDate = startDate.AddDate(0, 1, 0).Format(displayDateFormat)
	}

	// Convertir montos a palabras
	data.AmountToFinanceInWords, _ = p.numberConverter.ToInvoice(*financing.FinancingAmount, 2, currency)
	data.LotCostInWords, _ = p.numberConverter.ToInvoice(financing.Amount, 2, currency)

	if financing.DownPaymentBalance != nil {
		data.DownPaymentInWords, _ = p.numberConverter.ToInvoice(*financing.DownPaymentBalance, 2, currency)
	}

	return data, nil
}

func (p *PromissoryNotePDF) calculateAge(birthday string) (int, error) {
	birthDate, err := time.Parse(dateFormat, birthday)
	if err != nil {
		return 0, fmt.Errorf("formato de fecha de nacimiento inválido: %w", err)
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// Ajustar si no ha cumplido años este año
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age, nil
}

func (p *PromissoryNotePDF) getIdentityDocument(doc model.DocumentDomain) (string, error) {
	switch {
	case doc.DUI != nil:
		return *doc.DUI, nil
	case doc.NIT != nil:
		return *doc.NIT, nil
	case doc.Passport != nil:
		return *doc.Passport, nil
	default:
		return "", fmt.Errorf("el cliente no tiene documento de identidad")
	}
}

func (p *PromissoryNotePDF) validateRequiredFields(financing *model.FinancingDomain) error {
	if financing.TotalTerm == nil {
		return fmt.Errorf("el financiamiento no tiene plazo definido")
	}
	if financing.FinancingAmount == nil {
		return fmt.Errorf("el financiamiento no tiene monto definido")
	}
	if financing.MonthlyPayment == nil {
		return fmt.Errorf("el financiamiento no tiene pago mensual definido")
	}
	if financing.StartDate == nil {
		return fmt.Errorf("el financiamiento no tiene fecha de inicio")
	}
	return nil
}

// Optimizada función formatAmount usando strings.Builder
func formatAmount(amount float64) string {
	s := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := parts[1]

	if len(intPart) <= 3 {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(intPart) + (len(intPart)-1)/3 + 3) // Pre-allocate capacity

	// Insertar comas desde la derecha
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			builder.WriteByte(',')
		}
		builder.WriteRune(digit)
	}

	builder.WriteByte('.')
	builder.WriteString(decPart)

	return builder.String()
}

func NewPromissoryNotePDF(env *config.Env) port.PromissoryNoteService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	pagarePDF := pdf.NewGeneratePagarePDF()

	return &PromissoryNotePDF{
		api:               api,
		generatePagarePDF: pagarePDF,
		numberConverter:   numeroaletras.NewNumeroALetras(),
	}
}

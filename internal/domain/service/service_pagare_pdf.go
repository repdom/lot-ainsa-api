package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/financing"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"
	"strings"
	"time"
)

type PagarePDF struct {
	api               *financing.API
	generatePagarePDF pdf.GeneratePagarePDF
}

func (p PagarePDF) GenerateReport(financingId int) ([]byte, error) {
	loadFinancing, err := p.api.LoadFinancing("", "", "", financingId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var data pdf.PagareData

	fechaNac, err := time.Parse("2006-01-02", loadFinancing.Customer.Birthday)
	if err != nil {
		return nil, fmt.Errorf("error al parsear la fecha de nacimiento: %w", err)
	}

	fechaActual := time.Now()

	years := fechaActual.Year() - fechaNac.Year()
	if fechaNac.YearDay() > fechaActual.YearDay() {
		years--
	}

	data.NoteCode = fmt.Sprintf("%09f", loadFinancing.ID)
	data.ClientName = fmt.Sprintf("%s %s", loadFinancing.Customer.Names, loadFinancing.Customer.LastNames)
	data.IdentityDocument = loadFinancing.Customer.Document.DUI // loadFinancing.Customer.Document.Passport || loadFinancing.Customer.Document.NIT
	data.ClientAge = years
	data.Area = fmt.Sprintf("%.0f", loadFinancing.Lot.Area)
	data.LotNumber = loadFinancing.Lot.Number
	data.Address = loadFinancing.Customer.ResidentialAddress
	data.Block = loadFinancing.Lot.Polygon

	var termInYears = loadFinancing.TotalTerm / 12

	log.Printf("years term: %d", termInYears)

	data.TermInYears = termInYears
	data.InstallmentCount = loadFinancing.TotalTerm
	data.DateNow = fechaActual.Format("02/01/2006")
	data.InterestRate = fmt.Sprintf("%.2f", loadFinancing.InterestRate)
	data.AmountToFinance = formatAmount(loadFinancing.FinancingAmount)
	data.InstallmentAmount = formatAmount(loadFinancing.MonthlyPayment)
	data.DownPaymentPercentage = fmt.Sprintf("%.2f", loadFinancing.DownPaymentRate)
	data.DownPayment = formatAmount(loadFinancing.DownPaymentAmount)
	data.LotCost = formatAmount(loadFinancing.Lot.Price)
	data.Profession = loadFinancing.Customer.Financial.Occupation

	t, _ := time.Parse("2006-01-02", loadFinancing.StartDate)

	data.FirstPaymentDate = t.AddDate(0, 1, 0).Format("02/01/2006")
	data.AmountToFinanceInWords = ""
	data.LotCostInWords = ""
	data.DownPaymentInWords = ""

	return p.generatePagarePDF.GenerateReport(data)

}

func formatAmount(amount float64) string {
	// Convertir el número a string con 2 decimales
	s := fmt.Sprintf("%.2f", amount)

	// Separar parte entera y decimal
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := parts[1]

	// Insertar comas en la parte entera
	var result []string
	for i, r := range reverseString(intPart) {
		if i != 0 && i%3 == 0 {
			result = append(result, ",")
		}
		result = append(result, string(r))
	}

	// Revertir de nuevo la parte entera formateada
	formattedInt := reverseString(strings.Join(result, ""))

	return formattedInt + "." + decPart
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func NewPagarePDF(env *config.Env) port.PagareService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	api := financing.NewFinancingAPI(baseURL)
	pagarePDF := pdf.NewGeneratePagarePDF()
	return &PagarePDF{
		api:               api,
		generatePagarePDF: pagarePDF,
	}
}

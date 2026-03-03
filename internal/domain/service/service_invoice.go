package service

import (
	"be-lotsanmateo-api/internal/adapter/externalapi/client/company"
	downPayment "be-lotsanmateo-api/internal/adapter/externalapi/client/down/payment"
	"be-lotsanmateo-api/internal/adapter/externalapi/client/payment"
	"be-lotsanmateo-api/internal/adapter/externalapi/client/reservation"
	"be-lotsanmateo-api/internal/adapter/externalapi/model"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"
	"time"

	"github.com/user0608/numeroaletras"
)

const invoiceCurrency = "DOLARES"

type ServiceInvoice struct {
	numberConverter *numeroaletras.NumeroALetras
	payment         *payment.API
	downPayment     *downPayment.API
	reservation     *reservation.API
	company         *company.API
	invoice         *pdf.InvoicePagarePDF
}

func (s ServiceInvoice) InvoiceDownPayment(jwt, user, lang, name string, downPaymentId int) ([]byte, *string, error) {
	loadDownPayment, err := s.downPayment.GetDownPayment(jwt, user, lang, downPaymentId)
	if err != nil {
		return nil, nil, err
	}

	companyConfig, err := s.company.GetActiveCompanyConfiguration(jwt, user, lang)
	if err != nil {
		log.Printf("Error fetching company configuration: %v", err)
		// Si falla, podríamos usar valores por defecto o continuar con los que tenemos si fueran opcionales,
		// pero aquí los necesitamos para el reporte.
	}

	var clientName string
	var customer model.CustomerDomain
	var lot model.LotDomain
	var financing model.FinancingDomain

	if loadDownPayment.Financing != nil {
		financing = *loadDownPayment.Financing
		if financing.Customer != nil {
			customer = *financing.Customer
		}
		if financing.Lot != nil {
			lot = *financing.Lot
		}
	} else if loadDownPayment.Reservation != nil {
		customer = loadDownPayment.Reservation.Customer
		lot = loadDownPayment.Reservation.Lot
	}

	clientName = customer.Names + " " + customer.LastNames
	laCantidad, _ := s.numberConverter.ToInvoice(loadDownPayment.Amount, 2, invoiceCurrency)

	var document string
	var nit string

	if customer.Document.DUI == nil {
		if customer.Document.NIT != nil {
			nit = *customer.Document.NIT
		} else if customer.Document.Passport != nil {
			document = *customer.Document.Passport
		}
	} else {
		document = *customer.Document.DUI
	}

	var efectivo string
	if loadDownPayment.ReceiptNumber != nil && *loadDownPayment.ReceiptNumber != "" && *loadDownPayment.ReceiptNumber != "0" {
		efectivo = "SI"
	}

	var startingBalance string
	var remainingBalance string

	if loadDownPayment.Financing != nil && loadDownPayment.Financing.DownPaymentBalance != nil {
		startingBalance = fmt.Sprintf("%.2f", (loadDownPayment.Financing.Amount-*loadDownPayment.Financing.DownPaymentBalance)+loadDownPayment.Amount)
		remainingBalance = fmt.Sprintf("%.2f", loadDownPayment.Financing.Amount-*loadDownPayment.Financing.DownPaymentBalance)
	}

	data := pdf.ReciboPagoData{
		EmpresaNombre:           companyConfig.BusinessName,
		EmpresaDireccion:        companyConfig.Address,
		EmpresaTelefono:         companyConfig.Phone,
		EmpresaEmail:            companyConfig.Email,
		FechaEmision:            time.Now().Format("02/01/2006"),
		RecibidoPor:             name,
		CantidadPagadaEnNumeros: fmt.Sprintf("%.2f", loadDownPayment.Amount),
		RecibiDe:                clientName,
		LoteN:                   lot.Number,
		UbicacionDePoligono:     lot.Polygon,
		LaCantidadDe:            laCantidad,
		ConceptoDe:              "PAGO DE PRIMA",
		Direccion:               customer.ResidentialAddress,
		Dui:                     document,
		Nit:                     nit,
		ValorDelLote:            fmt.Sprintf("%.2f", lot.Price),
		Efectivo:                efectivo,
		AbonoACapital:           fmt.Sprintf("%.2f", loadDownPayment.Amount),
		InteresesMoratorios:     "0.00",
		InteresesNormales:       "0.00",
		SaldoAnterior:           startingBalance,
		SaldoActual:             remainingBalance,
	}

	var pdfData []byte
	pdfData, fail := s.invoice.GenerateReport(data)
	return pdfData, &clientName, fail
}

func (s ServiceInvoice) InvoicePayment(jwt, user, lang, name string, paymentId int) ([]byte, *string, error) {
	loadPayment, err := s.payment.GetPayment(jwt, user, lang, paymentId)
	if err != nil {
		return nil, nil, err
	}

	companyConfig, err := s.company.GetActiveCompanyConfiguration(jwt, user, lang)
	if err != nil {
		log.Printf("Error fetching company configuration: %v", err)
	}

	var clientName string
	var customer model.CustomerDomain
	var lot model.LotDomain

	if loadPayment.Financing != nil {
		if loadPayment.Financing.Customer != nil {
			customer = *loadPayment.Financing.Customer
		}
		if loadPayment.Financing.Lot != nil {
			lot = *loadPayment.Financing.Lot
		}
	} else if loadPayment.Reservation != nil {
		customer = loadPayment.Reservation.Customer
		lot = loadPayment.Reservation.Lot
	}

	clientName = customer.Names + " " + customer.LastNames
	laCantidad, _ := s.numberConverter.ToInvoice(loadPayment.Amount, 2, invoiceCurrency)

	var document string
	var nit string

	if customer.Document.DUI == nil {
		if customer.Document.NIT != nil {
			nit = *customer.Document.NIT
		} else if customer.Document.Passport != nil {
			document = *customer.Document.Passport
		}
	} else {
		document = *customer.Document.DUI
	}

	var efectivo string
	if loadPayment.ReceiptNumber != nil && *loadPayment.ReceiptNumber != "" && *loadPayment.ReceiptNumber != "0" {
		efectivo = "SI"
	}

	var principal string
	if loadPayment.Principal != nil {
		principal = fmt.Sprintf("%.2f", *loadPayment.Principal)
		log.Println(principal)
	}

	var startingBalance string
	if loadPayment.StartingBalance != nil {
		startingBalance = fmt.Sprintf("%.2f", *loadPayment.StartingBalance)
	}

	var remainingBalance string
	if loadPayment.RemainingBalance != nil {
		remainingBalance = fmt.Sprintf("%.2f", *loadPayment.RemainingBalance)
	}

	var interest string
	if loadPayment.Interest != nil {
		interest = fmt.Sprintf("%.2f", *loadPayment.Interest)
	}

	var penalty string
	if loadPayment.Penalty != nil {
		penalty = fmt.Sprintf("%.2f", *loadPayment.Penalty)
	}

	data := pdf.ReciboPagoData{
		EmpresaNombre:           companyConfig.BusinessName,
		EmpresaDireccion:        companyConfig.Address,
		EmpresaTelefono:         companyConfig.Phone,
		EmpresaEmail:            companyConfig.Email,
		FechaEmision:            time.Now().Format("02/01/2006"),
		RecibidoPor:             name,
		CantidadPagadaEnNumeros: fmt.Sprintf("%.2f", loadPayment.Amount),
		RecibiDe:                clientName,
		LoteN:                   lot.Number,
		UbicacionDePoligono:     lot.Polygon,
		LaCantidadDe:            laCantidad,
		ConceptoDe:              "PAGO Y/O ABONO A FINANCIAMIENTO",
		Direccion:               customer.ResidentialAddress,
		Dui:                     document,
		Nit:                     nit,
		ValorDelLote:            fmt.Sprintf("%.2f", lot.Price),
		Efectivo:                efectivo,
		SaldoAnterior:           startingBalance,
		SaldoActual:             remainingBalance,
		AbonoACapital:           principal,
		InteresesNormales:       interest,
		InteresesMoratorios:     penalty,
	}

	var pdfData []byte
	pdfData, fail := s.invoice.GenerateReport(data)
	return pdfData, &clientName, fail
}

func (s ServiceInvoice) InvoiceReservation(jwt, user, lang, name string, reservationId int) ([]byte, *string, error) {
	loadReservation, err := s.reservation.GetReservation(jwt, user, lang, reservationId)
	if err != nil {
		return nil, nil, err
	}

	companyConfig, err := s.company.GetActiveCompanyConfiguration(jwt, user, lang)
	if err != nil {
		log.Printf("Error fetching company configuration: %v", err)
	}

	clientName := loadReservation.Customer.Names + " " + loadReservation.Customer.LastNames
	laCantidad, _ := s.numberConverter.ToInvoice(loadReservation.Amount, 2, invoiceCurrency)

	var document string
	var nit string

	if loadReservation.Customer.Document.DUI == nil {
		if loadReservation.Customer.Document.NIT != nil {
			nit = *loadReservation.Customer.Document.NIT
		} else if loadReservation.Customer.Document.Passport != nil {
			document = *loadReservation.Customer.Document.Passport
		}
	} else {
		document = *loadReservation.Customer.Document.DUI
	}

	var efectivo string
	if loadReservation.ReceiptNumber != nil && *loadReservation.ReceiptNumber != "" && *loadReservation.ReceiptNumber != "0" {
		efectivo = "SI"
	}

	data := pdf.ReciboPagoData{
		EmpresaNombre:           companyConfig.BusinessName,
		EmpresaDireccion:        companyConfig.Address,
		EmpresaTelefono:         companyConfig.Phone,
		EmpresaEmail:            companyConfig.Email,
		FechaEmision:            time.Now().Format("02/01/2006"),
		RecibidoPor:             name,
		CantidadPagadaEnNumeros: fmt.Sprintf("%.2f", loadReservation.Amount),
		RecibiDe:                clientName,
		LoteN:                   loadReservation.Lot.Number,
		UbicacionDePoligono:     loadReservation.Lot.Polygon,
		LaCantidadDe:            laCantidad,
		ConceptoDe:              "RESERVA DE LOTE",
		Direccion:               loadReservation.Customer.ResidentialAddress,
		Dui:                     document,
		Nit:                     nit,
		ValorDelLote:            fmt.Sprintf("%.2f", loadReservation.Lot.Price),
		Efectivo:                efectivo,
		AbonoACapital:           fmt.Sprintf("%.2f", loadReservation.Amount),
		InteresesMoratorios:     "0.00",
		InteresesNormales:       "0.00",
		SaldoAnterior:           "0.00",
		SaldoActual:             fmt.Sprintf("%.2f", loadReservation.Lot.Price-loadReservation.Amount),
	}

	var pdfData []byte
	pdfData, fail := s.invoice.GenerateReport(data)
	return pdfData, &clientName, fail
}

func NewInvoiceService(env *config.Env) port.InvoiceService {
	baseURL := env.GetEnv("CUSTOMER_API_URL", "https://lot-db.rca-dev.com/")
	return &ServiceInvoice{
		numberConverter: numeroaletras.NewNumeroALetras(),
		payment:         payment.NewPaymentAPI(baseURL),
		downPayment:     downPayment.NewDownPaymentAPI(baseURL),
		reservation:     reservation.NewReservationAPI(baseURL),
		company:         company.NewCompanyAPI(baseURL),
		invoice:         pdf.NewInvoicePagarePDF(),
	}
}

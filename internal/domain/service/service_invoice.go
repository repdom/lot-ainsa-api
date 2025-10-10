package service

import (
	downPayment "be-lotsanmateo-api/internal/adapter/externalapi/client/down/payment"
	"be-lotsanmateo-api/internal/adapter/externalapi/client/payment"
	"be-lotsanmateo-api/internal/adapter/externalapi/client/reservation"
	"be-lotsanmateo-api/internal/adapter/report/pdf"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/internal/domain/port"
	"fmt"
	"log"

	"github.com/user0608/numeroaletras"
)

type ServiceInvoice struct {
	numberConverter *numeroaletras.NumeroALetras
	payment         *payment.API
	downPayment     *downPayment.API
	reservation     *reservation.API
	invoice         *pdf.InvoicePagarePDF
}

func (s ServiceInvoice) InvoiceDownPayment(jwt, user, lang, name string, downPaymentId int) ([]byte, *string, error) {
	loadDownPayment, err := s.downPayment.GetDownPayment(jwt, user, lang, downPaymentId)
	if err != nil {
		return nil, nil, err
	}

	clientName := loadDownPayment.Financing.Customer.Names + " " + loadDownPayment.Financing.Customer.LastNames
	laCantidad, _ := s.numberConverter.ToInvoice(loadDownPayment.Amount, 2, currency)

	var document string
	var nit string

	if loadDownPayment.Financing.Customer.Document.DUI == nil {
		if loadDownPayment.Financing.Customer.Document.NIT != nil {
			nit = *loadDownPayment.Financing.Customer.Document.NIT
		} else {
			document = *loadDownPayment.Financing.Customer.Document.Passport
		}
	} else {
		document = *loadDownPayment.Financing.Customer.Document.DUI
	}

	var efectivo string
	if loadDownPayment.ReceiptNumber != nil && *loadDownPayment.ReceiptNumber != "" && *loadDownPayment.ReceiptNumber != "0" {
		efectivo = "SI"
	}

	data := pdf.ReciboPagoData{
		RecibidoPor:             name,
		CantidadPagadaEnNumeros: fmt.Sprintf("%.2f", loadDownPayment.Amount),
		RecibiDe:                clientName,
		LoteN:                   loadDownPayment.Financing.Lot.Number,
		UbicacionDePoligono:     loadDownPayment.Financing.Lot.Polygon,
		LaCantidadDe:            laCantidad,
		ConceptoDe:              "PAGO DE PRIMA",
		Direccion:               loadDownPayment.Financing.Customer.ResidentialAddress,
		Dui:                     document,
		Nit:                     nit,
		ValorDelLote:            fmt.Sprintf("%.2f", loadDownPayment.Financing.Lot.Price),
		Efectivo:                efectivo,
		AbonoACapital:           fmt.Sprintf("%.2f", loadDownPayment.Amount),
		InteresesMoratorios:     "0.00",
		InteresesNormales:       "0.00",
		SaldoAnterior:           fmt.Sprintf("%.2f", (loadDownPayment.Financing.Amount-*loadDownPayment.Financing.DownPaymentBalance)+loadDownPayment.Amount),
		SaldoActual:             fmt.Sprintf("%.2f", loadDownPayment.Financing.Amount-*loadDownPayment.Financing.DownPaymentBalance),
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

	clientName := loadPayment.Financing.Customer.Names + " " + loadPayment.Financing.Customer.LastNames
	laCantidad, _ := s.numberConverter.ToInvoice(loadPayment.Amount, 2, currency)

	var document string
	var nit string

	if loadPayment.Financing.Customer.Document.DUI == nil {
		if loadPayment.Financing.Customer.Document.NIT != nil {
			nit = *loadPayment.Financing.Customer.Document.NIT
		} else {
			document = *loadPayment.Financing.Customer.Document.Passport
		}
	} else {
		document = *loadPayment.Financing.Customer.Document.DUI
	}

	var efectivo string
	if loadPayment.ReceiptNumber != nil && *loadPayment.ReceiptNumber != "" && *loadPayment.ReceiptNumber != "0" {
		efectivo = "SI"
	}

	log.Println(fmt.Sprintf("%.2f", *loadPayment.Principal))

	data := pdf.ReciboPagoData{
		RecibidoPor:             name,
		CantidadPagadaEnNumeros: fmt.Sprintf("%.2f", loadPayment.Amount),
		RecibiDe:                clientName,
		LoteN:                   loadPayment.Financing.Lot.Number,
		UbicacionDePoligono:     loadPayment.Financing.Lot.Polygon,
		LaCantidadDe:            laCantidad,
		ConceptoDe:              "PAGO Y/O ABONO A FINANCIAMIENTO",
		Direccion:               loadPayment.Financing.Customer.ResidentialAddress,
		Dui:                     document,
		Nit:                     nit,
		ValorDelLote:            fmt.Sprintf("%.2f", loadPayment.Financing.Lot.Price),
		Efectivo:                efectivo,
		SaldoAnterior:           fmt.Sprintf("%.2f", *loadPayment.StartingBalance),
		SaldoActual:             fmt.Sprintf("%.2f", *loadPayment.RemainingBalance),
		AbonoACapital:           fmt.Sprintf("%.2f", *loadPayment.Principal),
		InteresesNormales:       fmt.Sprintf("%.2f", *loadPayment.Interest),
		InteresesMoratorios:     fmt.Sprintf("%.2f", *loadPayment.Penalty),
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

	clientName := loadReservation.Customer.Names + " " + loadReservation.Customer.LastNames
	laCantidad, _ := s.numberConverter.ToInvoice(loadReservation.Amount, 2, currency)

	var document string
	var nit string

	if loadReservation.Customer.Document.DUI == nil {
		if loadReservation.Customer.Document.NIT != nil {
			nit = *loadReservation.Customer.Document.NIT
		} else {
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
		invoice:         pdf.NewInvoicePagarePDF(),
	}
}

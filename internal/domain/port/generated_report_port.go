package port

type PromissoryNoteService interface {
	GenerateReport(financingId int) ([]byte, error)
}

type ReportService interface {
	GenerateReport(financingId int) ([]byte, error)
}

type InvoiceService interface {
	InvoicePayment(paymentId int) ([]byte, error)
	InvoiceReservation(reservationId int) ([]byte, error)
	InvoicePremium(premiumId int) ([]byte, error)
}

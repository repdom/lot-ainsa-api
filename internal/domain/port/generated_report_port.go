package port

type PromissoryNoteService interface {
	GenerateReport(jwt, user, lang string, financingId int) ([]byte, *string, error)
}

type ReportService interface {
	GenerateReport(jwt, user, lang string, financingId int) ([]byte, *string, error)
}

type InvoiceService interface {
	InvoicePayment(jwt, user, lang, name string, paymentId int) ([]byte, *string, error)
	InvoiceDownPayment(jwt, user, lang, name string, downPaymentId int) ([]byte, *string, error)
	InvoiceReservation(jwt, user, lang, name string, reservationId int) ([]byte, *string, error)
}

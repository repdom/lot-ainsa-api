package pkg

import "be-lotsanmateo-api/internal/domain/model"

type PdfGenerator interface {
	Generate(data model.ResponseLoan) ([]byte, error)
}

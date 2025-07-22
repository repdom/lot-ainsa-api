package model

import "be-lotsanmateo-api/internal/adapter/externalapi/model/financing"

type PaymentLoanResponse struct {
	Interest      float64            `json:"interest"`
	Capital       float64            `json:"Capital"`
	Share         float64            `json:"share"`
	Penalty       float64            `json:"penalty"`
	AmountStart   float64            `json:"amountStart"`
	AmountBalance float64            `json:"AmountBalance"`
	Customer      financing.Customer `json:"customer"`
	Lot           financing.Lot      `json:"lot"`
}

package model

import "be-lotsanmateo-api/internal/adapter/externalapi/model"

type PaymentLoanResponse struct {
	Interest      float64               `json:"interest"`
	Capital       float64               `json:"Capital"`
	Share         float64               `json:"share"`
	Penalty       float64               `json:"interestOnArrearsAmount"`
	AmountStart   float64               `json:"amountStart"`
	AmountBalance float64               `json:"AmountBalance"`
	Customer      *model.CustomerDomain `json:"customer"`
	Lot           *model.LotDomain      `json:"lot"`
}

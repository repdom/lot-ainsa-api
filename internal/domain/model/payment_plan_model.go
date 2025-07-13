package model

type ResponseLoan struct {
	TotalAmount          float64          `json:"totalAmount"`
	NumberOfInstallments int              `json:"numberOfInstallments"`
	RateMonths           float64          `json:"rateMonths"`
	InterestsTotal       float64          `json:"interestsTotal"`
	MonthlyPayment       float64          `json:"monthlyPayments"`
	TotalPayments        float64          `json:"totalPayments"`
	PremiumRate          float64          `json:"premiumRate"`
	Premium              float64          `json:"premium"`
	FeeSimulation        []FreeSimulation `json:"feeSimulation"`
	Rate                 float64
	Years                float64
	Amount               float64
	DownPaymentRate      float64
}

type FreeSimulation struct {
	Amount       float64 `json:"amount"`
	Interest     float64 `json:"interest"`
	Capital      float64 `json:"capital"`
	BalanceStart float64 `json:"balanceStart"`
	BalanceLast  float64 `json:"balanceLast"`
	Payday       string  `json:"payday"`
}

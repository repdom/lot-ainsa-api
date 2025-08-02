package model

type RequestLoan struct {
	Rate    float64 `json:"rate"`
	Amount  float64 `json:"amount"`
	Months  int     `json:"months"`
	Payday  int     `json:"payday"`
	Premium float64 `json:"premium"`
}

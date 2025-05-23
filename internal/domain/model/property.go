package model

type Development struct {
	Name            string  `json:"name"`
	InterestRate    float64 `json:"interestRate"`
	DownPaymentRate float64 `json:"downPaymentRate"`
}

type Lots struct {
	Development Development `json:"development"`
	Price       float64     `json:"price"`
	Area        float64     `json:"size"`
	Number      int         `json:"number"`
	Polygon     string      `json:"polygon"`
}

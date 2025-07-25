package financing

import "time"

type Payment struct {
	ID               int       `json:"id"`
	Financing        string    `json:"financing"`
	PaymentDate      time.Time `json:"paymentDate"`
	Amount           float64   `json:"amount"`
	Principal        float64   `json:"principal"`
	Interest         float64   `json:"interest"`
	StartingBalance  float64   `json:"startingBalance"`
	RemainingBalance float64   `json:"remainingBalance"`
}

type Financings struct {
	ID                  float64       `json:"id"`
	Customer            Customer      `json:"customer"`
	Lot                 Lot           `json:"lot"`
	Reservation         Reservation   `json:"reservation"`
	Amount              float64       `json:"amount"`
	DownPaymentRate     float64       `json:"downPaymentRate"`
	DownPaymentAmount   float64       `json:"downPaymentAmount"`
	DownPaymentPending  float64       `json:"downPaymentPending"`
	DownPaymentBalance  float64       `json:"downPaymentBalance"`
	FinancingAmount     float64       `json:"financingAmount"`
	Balance             float64       `json:"balance"`
	InterestRate        float64       `json:"interestRate"`
	InterestRateMonthly float64       `json:"interestRateMonthly,omitempty"`
	TotalTerm           int           `json:"totalTerm,omitempty"`
	TermElapsed         int           `json:"termElapsed"`
	MissingTerm         int           `json:"missingTerm,omitempty"`
	MonthlyPayment      float64       `json:"monthlyPayment,omitempty"`
	Status              string        `json:"status"`
	StartDate           string        `json:"startDate,omitempty"`
	DownPayment         []DownPayment `json:"downPayment"`
	Payments            []Payment     `json:"payments"`
}

type Customer struct {
	ID                 int         `json:"id"`
	Names              string      `json:"names"`
	LastNames          string      `json:"lastNames"`
	Birthday           string      `json:"birthday"`
	Gender             Gender      `json:"gender"`
	CivilStatus        CivilStatus `json:"civilStatus"`
	Nationality        string      `json:"nationality"`
	ResidentialAddress string      `json:"residentialAddress"`
	City               string      `json:"city"`
	ZipCode            string      `json:"zipCode"`
	PhoneNumber        string      `json:"phoneNumber"`
	Email              string      `json:"email"`
	Document           Document    `json:"document"`
	PEP                PEP         `json:"pep"`
	Financial          Financial   `json:"financial"`
	Code               string      `json:"code"`
}

type Gender struct {
	ID     int    `json:"id"`
	Gender string `json:"gender"`
	Detail string `json:"detail"`
}

type CivilStatus struct {
	ID          int    `json:"id"`
	CivilStatus string `json:"civilStatus"`
}

type Document struct {
	ID       int     `json:"id"`
	DUI      string  `json:"dui"`
	NIT      *string `json:"nit"`
	Passport *string `json:"passport"`
}

type PEP struct {
	ID      int     `json:"id"`
	PEP     bool    `json:"pep"`
	Type    *string `json:"type"`
	Details *string `json:"details"`
}

type Financial struct {
	ID                   int    `json:"id"`
	Occupation           string `json:"occupation"`
	EmployerName         string `json:"employerName"`
	Position             string `json:"position"`
	IncomeSource         string `json:"incomeSource"`
	EstimatedIncomeRange string `json:"estimatedIncomeRange"`
	MainPurpose          string `json:"mainPurpose"`
}

type Lot struct {
	ID            int         `json:"id"`
	Development   Development `json:"development"`
	DevelopmentID *int        `json:"developmentId"`
	Number        string      `json:"number"`
	Polygon       string      `json:"polygon"`
	Area          float64     `json:"area"`
	AreaV2        float64     `json:"areaV2"`
	Price         float64     `json:"price"`
	PricePerV2    float64     `json:"pricePerV2"`
	Status        string      `json:"status"`
}

type Development struct {
	ID              int     `json:"id"`
	OwnerName       string  `json:"ownerName"`
	Location        string  `json:"location"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
	InterestRate    float64 `json:"interestRate"`
	DownPaymentRate float64 `json:"downPaymentRate"`
}

type Reservation struct {
	ID              int           `json:"id"`
	Lot             Lot           `json:"lot"`
	Customer        Customer      `json:"customer"`
	ReservationDate string        `json:"reservationDate"`
	ExpirationDate  string        `json:"expirationDate"`
	Status          string        `json:"status"`
	Amount          float64       `json:"amount"`
	DownPayment     []DownPayment `json:"downPayment"`
}

type DownPayment struct {
	ID                 int         `json:"id"`
	Reservation        Reservation `json:"reservation"`
	Amount             float64     `json:"amount"`
	DowPaymentStarting float64     `json:"dowPaymentStarting"`
	DownPaymentBalance float64     `json:"downPaymentBalance"`
	Reference          string      `json:"reference"`
}

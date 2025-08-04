package model

type Customer struct {
	ID                 int64       `json:"id,omitempty"`
	FirstName          string      `json:"names,omitempty"`
	LastName           string      `json:"lastNames,omitempty"`
	Birthday           string      `json:"birthday,omitempty"` // formato YYYY-MM-DD
	Gender             Gender      `json:"gender,omitempty"`
	CivilStatus        CivilStatus `json:"civilStatus,omitempty"`
	Nationality        string      `json:"nationality,omitempty"`
	ResidentialAddress string      `json:"residentialAddress,omitempty"`
	City               string      `json:"city,omitempty"`
	ZipCode            string      `json:"zipCode,omitempty"`
	PhoneNumber        string      `json:"phoneNumber,omitempty"`
	Email              string      `json:"email,omitempty"`
	Document           Document    `json:"document,omitempty"`
	Pep                Pep         `json:"pep,omitempty"`
	Financial          Financial   `json:"financial,omitempty"`
	Financings         []Financing `json:"financings,omitempty"`
}

type Gender struct {
	ID     int64  `json:"id,omitempty"`
	Gender string `json:"gender,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type CivilStatus struct {
	ID          int64  `json:"id,omitempty"`
	CivilStatus string `json:"civilStatus,omitempty"`
}

type Document struct {
	ID       int64  `json:"id,omitempty"`
	DUI      string `json:"dui,omitempty"`
	NIT      string `json:"nit,omitempty"`
	Passport string `json:"passport,omitempty"`
}

type Pep struct {
	ID      int64  `json:"id,omitempty"`
	Pep     bool   `json:"pep,omitempty"`
	Type    string `json:"type,omitempty"`
	Details string `json:"details,omitempty"`
}

type Financial struct {
	ID                   int64  `json:"id,omitempty"`
	Occupation           string `json:"occupation,omitempty"`
	EmployerName         string `json:"employerName,omitempty"`
	Position             string `json:"position,omitempty"`
	IncomeSource         string `json:"incomeSource,omitempty"`
	EstimatedIncomeRange string `json:"estimatedIncomeRange,omitempty"`
	MainPurpose          string `json:"mainPurpose,omitempty"`
	Customer             string `json:"customer,omitempty"`
}

type Financing struct {
	ID                  int64       `json:"id,omitempty"`
	Customer            string      `json:"customer,omitempty"`
	Lot                 Lot         `json:"lot,omitempty"`
	Reservation         Reservation `json:"reservation,omitempty"`
	Amount              float64     `json:"amount,omitempty"`
	DownPaymentRate     float64     `json:"downPaymentRate,omitempty"`
	DownPaymentAmount   float64     `json:"downPaymentAmount,omitempty"`
	DownPaymentPending  float64     `json:"downPaymentPending,omitempty"`
	DownPaymentBalance  float64     `json:"downPaymentBalance,omitempty"`
	FinancingAmount     float64     `json:"financingAmount,omitempty"`
	Balance             float64     `json:"balance,omitempty"`
	InterestRate        float64     `json:"interestRate,omitempty"`
	InterestRateMonthly float64     `json:"interestRateMonthly,omitempty"`
	TotalTerm           int         `json:"totalTerm,omitempty"`
	TermElapsed         int         `json:"termElapsed,omitempty"`
	MissingTerm         int         `json:"missingTerm,omitempty"`
	MonthlyPayment      float64     `json:"monthlyPayment,omitempty"`
	Status              string      `json:"status,omitempty"`
	StartDate           string      `json:"startDate,omitempty"` // formato YYYY-MM-DD
}

type Reservation struct {
	ID              int64   `json:"id,omitempty"`
	Lot             Lot     `json:"lot,omitempty"`
	Customer        string  `json:"customer,omitempty"`
	ReservationDate string  `json:"reservationDate,omitempty"`
	ExpirationDate  string  `json:"expirationDate,omitempty"`
	Status          string  `json:"status,omitempty"`
	Amount          float64 `json:"amount,omitempty"`
}

type Lot struct {
	ID          int64       `json:"id,omitempty"`
	Development Development `json:"development,omitempty"`
	Number      string      `json:"number,omitempty"`
	Polygon     string      `json:"polygon,omitempty"`
	Area        float64     `json:"area,omitempty"`
	Price       float64     `json:"price,omitempty"`
	Status      string      `json:"status,omitempty"`
}

type Development struct {
	ID              int64    `json:"id,omitempty"`
	OwnerName       string   `json:"ownerName,omitempty"`
	Location        string   `json:"location,omitempty"`
	Description     string   `json:"description,omitempty"`
	Status          string   `json:"status,omitempty"`
	InterestRate    float64  `json:"interestRate,omitempty"`
	DownPaymentRate float64  `json:"downPaymentRate,omitempty"`
	Lots            []string `json:"lots,omitempty"`
}

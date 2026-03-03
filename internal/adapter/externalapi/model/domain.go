package model

type PaymentDomain struct {
	ID                 int                `json:"id"`
	Reservation        *ReservationDomain `json:"reservation,omitempty"`
	Financing          *FinancingDomain   `json:"financing,omitempty"`
	PaymentDate        *string            `json:"paymentDate,omitempty"`
	Amount             float64            `json:"amount"`
	Principal          *float64           `json:"principal,omitempty"`
	Interest           *float64           `json:"interest,omitempty"`
	Penalty            *float64           `json:"penalty,omitempty"`
	StartingBalance    *float64           `json:"startingBalance,omitempty"`
	RemainingBalance   *float64           `json:"remainingBalance,omitempty"`
	ReceiptNumber      *string            `json:"receiptNumber"`
	ReferenceNumber    *string            `json:"referenceNumber,omitempty"`
	DowPaymentStarting *float64           `json:"dowPaymentStarting,omitempty"`
	DownPaymentBalance *float64           `json:"downPaymentBalance,omitempty"`
}

type FinancingDomain struct {
	ID                     int                  `json:"id"`
	Customer               *CustomerDomain      `json:"customer,omitempty"`
	Lot                    *LotDomain           `json:"lot,omitempty"`
	Reservation            *ReservationDomain   `json:"reservation,omitempty"`
	Amount                 float64              `json:"amount"`
	DownPaymentRate        *float64             `json:"downPaymentRate,omitempty"`
	DownPaymentAmount      *float64             `json:"downPaymentAmount,omitempty"`
	DownPaymentPending     *float64             `json:"downPaymentPending,omitempty"`
	DownPaymentBalance     *float64             `json:"downPaymentBalance,omitempty"`
	FinancingAmount        *float64             `json:"financingAmount,omitempty"`
	FinancingBalance       float64              `json:"financingBalance"`
	FinancingAmountPending *float64             `json:"financingAmountPending,omitempty"`
	InterestBalance        *float64             `json:"interestBalance,omitempty"`
	InterestRate           *float64             `json:"interestRate,omitempty"`
	InterestRateMonthly    *float64             `json:"interestRateMonthly,omitempty"`
	TotalTerm              *int                 `json:"totalTerm,omitempty"`
	TermElapsed            int                  `json:"termElapsed"`
	MissingTerm            *int                 `json:"missingTerm,omitempty"`
	MonthlyPayment         *float64             `json:"monthlyPayment,omitempty"`
	Status                 string               `json:"status"`
	StartDate              *string              `json:"startDate,omitempty"`
	DownPayment            *[]DownPaymentDomain `json:"downPayment,omitempty"`
	PaymentPlans           *[]interface{}       `json:"paymentPlans,omitempty"` // ignorado de momento
	Payments               *[]PaymentDomain     `json:"payments,omitempty"`
}

type CustomerDomain struct {
	ID                 *int               `json:"id,omitempty"`
	Names              string             `json:"names"`
	LastNames          string             `json:"lastNames"`
	Birthday           string             `json:"birthday"`
	Gender             *GenderDomain      `json:"gender,omitempty"`
	CivilStatus        *CivilStatusDomain `json:"civilStatus,omitempty"`
	Nationality        *string            `json:"nationality,omitempty"`
	ResidentialAddress string             `json:"residentialAddress"`
	City               *string            `json:"city,omitempty"`
	ZipCode            *string            `json:"zipCode,omitempty"`
	PhoneNumber        string             `json:"phoneNumber"`
	Email              *string            `json:"email,omitempty"`
	Document           DocumentDomain     `json:"document"`
	Pep                *PepDomain         `json:"pep,omitempty"`
	Financial          *FinancialDomain   `json:"financial,omitempty"`
	Code               *string            `json:"code,omitempty"`
	Profession         *string            `json:"profession,omitempty"`
}

type GenderDomain struct {
	ID     *int    `json:"id,omitempty"`
	Gender string  `json:"gender"`
	Detail *string `json:"detail,omitempty"`
}

type CivilStatusDomain struct {
	ID          *int   `json:"id,omitempty"`
	CivilStatus string `json:"civilStatus"`
}

type DocumentDomain struct {
	ID       *int    `json:"id,omitempty"`
	DUI      *string `json:"dui,omitempty"`
	NIT      *string `json:"nit,omitempty"`
	Passport *string `json:"passport,omitempty"`
}

type PepDomain struct {
	ID           *int    `json:"id"`
	Pep          bool    `json:"pep"`
	Type         *string `json:"type"`
	Details      *string `json:"details"`
	FullName     *string `json:"fullName"`
	Title        *string `json:"title"`
	Relationship *string `json:"relationship"`
}

type FinancialDomain struct {
	ID                   *int   `json:"id,omitempty"`
	Occupation           string `json:"occupation"`
	EmployerName         string `json:"employerName"`
	Position             string `json:"position"`
	IncomeSource         string `json:"incomeSource"`
	EstimatedIncomeRange string `json:"estimatedIncomeRange"`
	MainPurpose          string `json:"mainPurpose"`
}

type LotDomain struct {
	ID            int                 `json:"id"`
	Development   DevelopmentDomain   `json:"development"`
	DevelopmentId *int                `json:"developmentId,omitempty"`
	Number        string              `json:"number"`
	Polygon       string              `json:"polygon"`
	Area          float64             `json:"area"`
	AreaV2        float64             `json:"areaV2"`
	Price         float64             `json:"price"`
	PricePerV2    float64             `json:"pricePerV2"`
	Status        string              `json:"status"`
	Financings    []FinancingDomain   `json:"financings,omitempty"`
	Reservations  []ReservationDomain `json:"reservations,omitempty"`
}

type DevelopmentDomain struct {
	ID                int         `json:"id"`
	OwnerName         string      `json:"ownerName"`
	Name              string      `json:"name"`
	Location          string      `json:"location"`
	Description       string      `json:"description"`
	Status            string      `json:"status"`
	InterestRate      float64     `json:"interestRate"`
	InterestOnArrears float64     `json:"interestOnArrears"`
	DownPaymentRate   float64     `json:"downPaymentRate"`
	Lots              interface{} `json:"lots,omitempty"`
}

type ReservationDomain struct {
	ID              int                 `json:"id"`
	Lot             LotDomain           `json:"lot"`
	Customer        CustomerDomain      `json:"customer"`
	ReservationDate string              `json:"reservationDate"`
	ExpirationDate  string              `json:"expirationDate"`
	Status          string              `json:"status"`
	Amount          float64             `json:"amount"`
	ReceiptNumber   *string             `json:"receiptNumber"`
	ReferenceNumber string              `json:"referenceNumber"`
	DownPayment     []DownPaymentDomain `json:"downPayment,omitempty"`
	Financings      []FinancingDomain   `json:"financings,omitempty"`
}

type DownPaymentDomain struct {
	ID                 int                `json:"id"`
	Reservation        *ReservationDomain `json:"reservation,omitempty"`
	Financing          *FinancingDomain   `json:"financing,omitempty"`
	Amount             float64            `json:"amount"`
	DowPaymentStarting *float64           `json:"dowPaymentStarting"`
	DownPaymentBalance float64            `json:"downPaymentBalance"`
	ReceiptNumber      *string            `json:"receiptNumber"`
	ReferenceNumber    *string            `json:"referenceNumber,omitempty"`
}

type CompanyConfigurationDomain struct {
	ID           int    `json:"id"`
	RNC          string `json:"rnc"`
	BusinessName string `json:"businessName"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Mobile       string `json:"mobile"`
	Whatsapp     string `json:"whatsapp"`
	Active       bool   `json:"active"`
}

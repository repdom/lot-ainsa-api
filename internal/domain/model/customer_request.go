package model

type RequestCustomerOnboarding struct {
	Address              string `json:"address"`
	BirthDate            string `json:"birthDate"` // Formato: YYYY-MM-DD
	City                 string `json:"city"`
	DetailsPep           string `json:"detailsPep"`
	DocumentNumber       string `json:"documentNumber"`
	DocumentType         string `json:"documentType"`
	Email                string `json:"email"`
	Employer             string `json:"employer"`
	Gender               string `json:"gender"`
	IsExposedPolitically string `json:"isExposedPolitically"`
	Lastnames            string `json:"lastnames"`
	MaritalStatus        string `json:"maritalStatus"`
	Names                string `json:"names"`
	Nationality          string `json:"nationality"`
	Occupation           string `json:"occupation"`
	Phone                string `json:"phone"`
	Position             string `json:"position"`
	PostalCode           string `json:"postalCode"`
	RangeIncome          string `json:"rangeIncome"`
	RelationFinancial    string `json:"relationFinancial"`
	SourceOfIncome       string `json:"sourceOfIncome"`
	FullName             string `json:"fullName"`
	Title                string `json:"title"`
	Relationship         string `json:"relationship"`
	ProfessionDUI        string `json:"professionDUI"`
}

package models

type ExternalCountry struct {
	Name       string             `json:"name"`
	Capital    string             `json:"capital"`
	Region     string             `json:"region"`
	Population int64              `json:"population"`
	Flag       string             `json:"flag"`
	Currencies []ExternalCurrency `json:"currencies"`
}

type ExternalCurrency struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

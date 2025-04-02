package models

// RegistrationRequest defines the structure for user registration data
type RegistrationRequest struct {
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Country          string  `json:"country"`
	Address          string  `json:"address"`
	PostalCode       string  `json:"postalCode"`
	DogeAddress      string  `json:"dogeAddress"`
	Size             string  `json:"size"`
	BName            string  `json:"bname"`
	BEmail           string  `json:"bemail"`
	BCountry         string  `json:"bcountry"`
	BAddress         string  `json:"baddress"`
	BPostalCode      string  `json:"bpostalCode"`
	Amount           float64 `json:"amount"`
	Sku              string  `json:"sku"`
	PaytoDogeAddress string  `json:"paytoDogeAddress,omitempty"`
}

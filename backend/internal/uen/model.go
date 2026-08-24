package uen

type ValidateRequest struct {
	UEN string `json:"uen"`
}

type ValidateResponse struct {
	UEN     string `json:"uen"`
	Valid   bool   `json:"valid"`
	Format  string `json:"format,omitempty"`
	Message string `json:"message"`
}

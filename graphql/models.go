package main
type Account struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Orders []Order `json:"orders"`
}

type AccountResponse struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email,omitempty"`
	Orders []Order `json:"orders,omitempty"`
}

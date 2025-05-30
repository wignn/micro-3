package model


type Account struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Password string `json:"password,omitempty"`
	Email string `json:"email,omitempty"`
}

type AccountResponse struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email,omitempty"`
}

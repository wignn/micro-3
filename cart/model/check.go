package model

type CartProduct struct {
	ProductID string `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
}

type Cart struct {
	ID        string        `json:"id"`
	AccountID string        `json:"account_id"`
	CreatedAt string        `json:"created_at"`
	Products  []CartProduct `json:"products"`
}

type CartPutRequest struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	CreatedAt string `json:"created_at"`
	ProductID string `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
}

type CartPutResponse struct {
	ID        string `json:"id"`
	Quantity  uint32 `json:"quantity"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Description string `json:"description"`
}

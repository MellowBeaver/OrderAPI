// The structures of database tables.
package model

//Creating Order Item struct

type OrderItem struct {
	Id    int     `json:"id"`
	Desc  string  `json:"description"`
	Price float32 `json:"price"`
	Qty   int     `json:"quantity"`
}

//Creating Order struct

type Order struct {
	Id           string    `json:"id"`
	Status       string    `json:"status"`
	Items        OrderItem `json:"items"`
	Total        float32   `json:"total"`
	CurrencyUnit string    `json:"currencyUnit"`
}

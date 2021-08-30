package engine

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// Order Type
type Order struct {
	ID        string          `json:"id"`
	Side      Side            `json:"side"`
	Quantity  decimal.Decimal `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
	Timestamp int64           `json:"timestamp"`
}

// Convert order to struct from json
func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

// Convert order to json from order struct
func (order *Order) toJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}

package engine

import "encoding/json"

type Trade struct {
	TakerOrderID string `json:"taker_order_id"`
	MakerOrderID string `json:"maker_order_id"`
	Quantity     uint64 `json:"quantity"`
	Price        uint64 `json:"price"`
	Timestamp    int64  `json:"timestamp"`
}

// struct to json
func (trade *Trade) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}

func (trade *Trade) ToJSON() []byte {
	str, _ := json.Marshal(trade)
	return str
}

package engine

import (
	"encoding/json"
	"reflect"
)

type Side int

// 0 = sell 1 = buy
const (
	Sell Side = iota
	Buy
)

func (s Side) String() string {
	if s == Buy {
		return "buy"
	}
	return "sell"
}

// implement Marshalling and Unmarshalling
// will convert struct to json
func (s Side) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

// will convert json to struct
func (s *Side) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"buy"`:
		*s = Buy
	case `"sell"`:
		*s = Sell
	default:
		return &json.UnsupportedValueError{
			Value: reflect.New(reflect.TypeOf(data)),
			Str:   string(data),
		}
	}
	return nil
}

package bitmex

import (
	"encoding/json"
)

type Response struct {
	Success   bool          `json:"success,omitempty"`
	Subscribe string        `json:"subscribe,omitempty"`
	Request   interface{}   `json:"request,omitempty"`
	Table     string        `json:"table,omitempty"`
	Action    string        `json:"action,omitempty"`
	Data      []interface{} `json:"data,omitempty"`
}

// DecodeMessage decodes the raw JSON from a BitMEX API response
// and returns a Response
func DecodeMessage(message []byte) (Response, error) {
	var res Response
	err := json.Unmarshal(message, &res)

	return res, err
}

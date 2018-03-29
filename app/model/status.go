package model

import "encoding/json"

// Status represents a status response message
type Status struct {
	Code    int    `json:"code"`
	Text    string `json:"text"`
	Network int    `json:"network,omitempty"`
}

// String converts object to string
func (s Status) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

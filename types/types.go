package types

import (
	"bytes"
	"encoding/json"
)

type ObjectMap = map[string]interface{}

func Encode(om ObjectMap) *bytes.Buffer {
	var buff bytes.Buffer
	_ = json.NewEncoder(&buff).Encode(om)

	return &buff
}

type SearchResponse struct {
	Hits struct {
		Hist []json.RawMessage `json:"hits"`
	} `json:"hits"`
}

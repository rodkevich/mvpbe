package v1

import (
	"encoding/json"
	"fmt"
	"io"
)

type SampleItemRequest struct {
	ID     int    `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

func (i *SampleItemRequest) Bind(body io.ReadCloser) error {
	defer body.Close()
	err := json.NewDecoder(body).Decode(i)
	if err != nil {
		return fmt.Errorf("data bind: %w", err)
	}
	return nil
}

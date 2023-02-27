package model

import (
	"time"
)

var (
	ItemCreated       = "CREATED"
	ItemPending       = "PENDING"
	ItemComplete      = "COMPLETE"
	ItemDeletePending = "DEL_PEND"
	ItemDeleted       = "DELETED"
)

type SampleItem struct {
	ID         int       `json:"id,omitempty,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	FinishTime time.Time `json:"finish_time,omitempty"`
	Status     string    `json:"status,omitempty"`
}

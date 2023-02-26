package model

import "time"

var (
	ItemCreated       = "CREATED"
	ItemPending       = "PENDING"
	ItemComplete      = "COMPLETE"
	ItemDeletePending = "DEL_PEND"
	ItemDeleted       = "DELETED"
)

type SampleItem struct {
	ID         int
	StartTime  time.Time
	FinishTime time.Time
	Status     string
}

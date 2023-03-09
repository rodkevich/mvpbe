package model

import (
	"time"

	v1 "github.com/rodkevich/mvpbe/pkg/api/v1"
)

const undoRedoTimeout = 1 * time.Minute

var (
	ItemCreated  = "CREATED"
	ItemPending  = "PENDING"
	ItemComplete = "COMPLETE"
	ItemDeleted  = "DELETED"
)

type SampleItem struct {
	ID         int       `json:"id,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	FinishTime time.Time `json:"finish_time,omitempty"`
	Status     string    `json:"status,omitempty"`
}

type SomeProcessingTask struct {
	SampleItem `json:"sample_item"`
	TraceID    string `json:"trace_id"`
}

func (c *SampleItem) Expired() bool {
	return c.FinishTime.Add(undoRedoTimeout).Unix() < v1.TimeNow().Unix()
}

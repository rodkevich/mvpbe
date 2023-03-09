package itemsprocessor

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rodkevich/mvpbe/internal/domain/itemsprocessor/datasource"
	"github.com/rodkevich/mvpbe/internal/domain/itemsprocessor/model"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

// Processor represents sample usage of sample domain
type Processor interface {
	Readiness() error
	UpdateItem(ctx context.Context, sip *model.SomeProcessingTask) error
}

// Items implements Processor
type Items struct {
	db  *datasource.SampleProcessorDB
	rmq rabbitmq.AMQPPublisher
}

// UpdateItem ...
func (i *Items) UpdateItem(ctx context.Context, proc *model.SomeProcessingTask) error {
	proc.FinishTime = api.TimeNow()

	err := i.db.UpdateStatusExampleTrx(ctx, &proc.SampleItem)
	if err != nil {
		return fmt.Errorf("remote update failed: %w", err)
	}

	dataBytes, err := json.Marshal(proc.SampleItem)
	if err != nil {
		return fmt.Errorf("json.marshal failed: %w", err)
	}

	// if item was deleted from processing publish it somewhere in another que
	if proc.Status == model.ItemDeleted {
		return i.rmq.Publish(
			ctx, exExchangeName, exBindingKeyItemsReadiness,
			amqp.Publishing{
				Headers: map[string]interface{}{
					"example-item-trace-id": fmt.Sprintf(api.ExampleItemsTracingFormat, proc.ID),
				},
				Timestamp: api.TimeNow(),
				Body:      dataBytes,
			})
	}

	// publish to workers que
	return i.rmq.Publish(
		ctx, exExchangeName, exBindingKeyItemsProcessing,
		amqp.Publishing{
			Headers: map[string]interface{}{
				"example-item-trace-id": fmt.Sprintf(api.ExampleItemsTracingFormat, proc.ID),
			},
			Timestamp: api.TimeNow(),
			Body:      dataBytes,
		})
}

// Readiness of domain
func (i *Items) Readiness() error {
	return i.db.Readiness()
}

// NewItemsDomain constructor
func NewItemsDomain(repo *datasource.SampleProcessorDB, pbl rabbitmq.AMQPPublisher) *Items {
	itemsUsage := &Items{db: repo, rmq: pbl}
	return itemsUsage
}

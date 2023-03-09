package itemsprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
			ctx, exExchangeNameItems, exBindingKeyItemsReadiness,
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
		ctx, exExchangeNameItems, exBindingKeyItems,
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
func NewItemsDomain(ctx context.Context, repo *datasource.SampleProcessorDB, pbl rabbitmq.AMQPPublisher) *Items {
	channel := pbl.GetChannel()
	configureExchanges(channel)
	itemsUsage := &Items{db: repo, rmq: pbl}

	itemsCh, err := channel.Consume(exQueueNameItems, exConsumerNameItems, false, false, false, false, nil)
	if err != nil {
		log.Fatal("err := channel.Consume")
	}

	go func() {
		runExampleItemsConsumer(ctx, itemsUsage, itemsCh)
	}()

	return itemsUsage
}

func configureExchanges(ch *amqp.Channel) {
	log.Println("configuring rabbit ")
	err := ch.ExchangeDeclare(exExchangeNameItems, exExchangeKindItems, true, false, false, false, nil)
	if err != nil {
		log.Fatal("err := ch.ExchangeDeclare: ", err)
	}

	// configure some ques and their bindings
	for k, v := range map[string]string{
		exQueueNameItems:   exBindingKeyItems,
		exQueueNameResults: exBindingKeyItemsReadiness,
	} {
		q, err := ch.QueueDeclare(k, true, false, false, false, nil)
		if err != nil {
			log.Fatal("err := ch.QueueDeclare: ", err)
		}

		err = ch.QueueBind(q.Name, v, exExchangeNameItems, false, nil)
		if err != nil {
			log.Fatal("err := ch.QueueBind: ", err)
		}
	}
}

package items_processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rodkevich/mvpbe/internal/domain/items_processor/datasource"
	"github.com/rodkevich/mvpbe/internal/domain/items_processor/model"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

// ItemsSampleProcessUsage represents sample usage of sample domain
type ItemsSampleProcessUsage interface {
	Readiness() error
	UpdateItem(ctx context.Context, m *model.SampleItem) error
}

// Items implements ItemsSampleProcessUsage
type Items struct {
	db  *datasource.SampleProcessorDB
	rmq rabbitmq.AMQPPublisher
}

// UpdateItem ...
func (i *Items) UpdateItem(ctx context.Context, m *model.SampleItem) error {
	m.FinishTime = api.TimeNow()
	err := i.db.UpdateStatusExampleTrx(ctx, m)
	if err != nil {
		return fmt.Errorf("remote update failed: %w", err)
	}

	dataBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("json.marshal failed: %w", err)
	}

	return i.rmq.Publish(
		ctx, exampleItemsExchangeName, exampleItemsBindingKey,
		amqp.Publishing{
			Headers:   map[string]interface{}{"example-item-trace-id": m.ID},
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

	itemsCh, err := channel.Consume(exampleItemsQueueName, exampleItemsConsumerName, false, false, false, false, nil)
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
	err := ch.ExchangeDeclare(exampleItemsExchangeName, exampleItemsExchangeKind, true, false, false, false, nil)
	if err != nil {
		log.Fatal("err := ch.ExchangeDeclare")
	}
	queue, err := ch.QueueDeclare(exampleItemsQueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("err := ch.QueueDeclare")
	}
	err = ch.QueueBind(queue.Name, exampleItemsBindingKey, exampleItemsExchangeName, false, nil)
	if err != nil {
		log.Fatal("err := ch.QueueBind")
	}
}

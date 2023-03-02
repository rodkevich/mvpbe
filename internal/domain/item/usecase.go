package item

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"

	"github.com/rodkevich/mvpbe/internal/domain/item/datasource"
	"github.com/rodkevich/mvpbe/internal/domain/item/model"
	"github.com/rodkevich/mvpbe/pkg/rabbitmq"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

//go:generate mockery --name ItemsSampleUsage --case underscore  --output mocks/

// ItemsSampleUsage represents sample usage of sample domain
type ItemsSampleUsage interface {
	Readiness() error
	AllDatabases(ctx context.Context) ([]string, error)
	AddItem(ctx context.Context, m *model.SampleItem) error
	UpdateItem(ctx context.Context, m *model.SampleItem) error
	GetItem(ctx context.Context, id int) (*model.SampleItem, error)
	ListItems(ctx context.Context) ([]*model.SampleItem, error)
}

// Items implements ItemsSampleUsage
type Items struct {
	db  *datasource.SampleDB
	rmq rabbitmq.AMQPPublisher
}

// AddItem ...
func (i *Items) AddItem(ctx context.Context, m *model.SampleItem) error {
	m.StartTime = api.TimeNow()
	m.FinishTime = api.TimeNow()
	m.Status = model.ItemCreated

	err := i.db.AddItemExampleTrx(ctx, m)
	if err != nil {
		return fmt.Errorf("remote add failed: %w", err)
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

// UpdateItem ...
func (i *Items) UpdateItem(ctx context.Context, m *model.SampleItem) error {
	m.FinishTime = api.TimeNow()
	return i.db.UpdateStatusExampleTrx(ctx, m)
}

// GetItem ...
func (i *Items) GetItem(ctx context.Context, id int) (*model.SampleItem, error) {
	example, err := i.db.GetItemExample(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("remote get failed: %w", err)
	}

	if example.Status == model.ItemDeleted {
		return nil, errDeletedItem
	}
	return example, nil
}

// ListItems ...
func (i *Items) ListItems(ctx context.Context) ([]*model.SampleItem, error) {
	// TODO implement me
	panic("implement me")
}

// Readiness of domain
func (i *Items) Readiness() error {
	return i.db.Readiness()
}

// AllDatabases sample method to get with all db names
func (i *Items) AllDatabases(ctx context.Context) ([]string, error) {
	return i.db.AllDatabases(ctx)
}

// NewItemsDomain constructor
func NewItemsDomain(ctx context.Context, repo *datasource.SampleDB, pbl rabbitmq.AMQPPublisher) *Items {
	channel := pbl.GetChannel()
	configureExchanges(channel)
	itemsUsage := &Items{db: repo, rmq: pbl}

	itemsCh, err := channel.Consume(
		exampleItemsQueueName,
		"items-consumer-0001",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("err := channel.Consume")
	}

	go func() {
		eg, ctx := errgroup.WithContext(ctx)
		for i := 0; i <= exampleItemsAMQPConcurrency; i++ {
			eg.Go(func(ctx context.Context, itemsCh <-chan amqp.Delivery, workerID int) func() error {
				log.Printf("starting consumer id: %d, for items queue: %s", workerID, exampleItemsQueueName)

				return func() error {
					for {
						select {
						case <-ctx.Done():
							log.Printf("items consumer ctx done: %v", ctx.Err())
							return ctx.Err()

						case msg, ok := <-itemsCh:
							if !ok {
								log.Printf("NOT OK items channel closed for queue: %s", exampleItemsQueueName)
								return errors.New("items channel closed")
							}
							log.Printf("Items consumer: id: %d, msg data: %s, headers: %+v", workerID, string(msg.Body), msg.Headers)

							m := model.SampleItem{}
							err := json.Unmarshal(msg.Body, &m)
							if err != nil {
								return err
							}
							m.Status = model.ItemPending
							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								return err
							}
							mod, err := itemsUsage.GetItem(ctx, m.ID)
							if err != nil {
								return err
							}
							log.Printf("FINAL item view: %+v\n", mod)

							log.Printf("Consumer [ACK] delivery: id: %d, msg data: %s, headers: %+v", workerID, string(msg.Body), msg.Headers)
							_ = msg.Ack(false)
						}
					}
				}
			}(ctx, itemsCh, i))
		}
		_ = eg.Wait()
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

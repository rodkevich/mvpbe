package item

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"

	"github.com/rodkevich/mvpbe/internal/domain/item/model"
)

func runExampleItemsConsumer(ctx context.Context, itemsUsage *Items, itemsCh <-chan amqp.Delivery) {
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

						// TODO remove
						const duration = 3 * time.Second

						switch m.Status {
						case model.ItemCreated:
							time.Sleep(duration)
							pending := model.ItemPending
							m.Status = pending
							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								log.Println("items consumer got error with model.ItemCreated: itemsUsage.UpdateItem: ", err)
								continue
							}
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error with case model.ItemCreated: msg.Ack(false): ", err)
								continue
							}

						case model.ItemPending:
							time.Sleep(duration)
							pending := model.ItemComplete
							m.Status = pending
							err = itemsUsage.UpdateItem(ctx, &m)
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error with case model.ItemPending: msg.Ack(false): ", err)
								continue
							}

						case model.ItemComplete:
							time.Sleep(duration)
							pending := model.ItemDeleted
							m.Status = pending
							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								log.Println("items consumer got error with model.ItemComplete: itemsUsage.UpdateItem: ", err)
								continue
							}
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error with case model.ItemComplete: msg.Ack(false): ", err)
								continue
							}
							continue

						case model.ItemDeleted:
							err := msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error with case model.ItemDeleted: msg.Ack(false): ", err)
								continue
							}
							continue
						default:
							continue
						}
					}
				}
			}
		}(ctx, itemsCh, i))
	}
	_ = eg.Wait()
}

package item

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"

	"github.com/rodkevich/mvpbe/internal/domain/item/model"
)

func runExampleItemsConsumer(ctx context.Context, itemsUsage ItemsSampleUsage, itemsCh <-chan amqp.Delivery) {
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
							log.Printf("error: items channel closed for queue: %s", exampleItemsQueueName)
							return errors.New("items channel closed")
						}
						log.Printf("Items consumer: id: %d, msg data: %s, headers: %+v", workerID, string(msg.Body), msg.Headers)

						m := model.SampleItem{}
						err := json.Unmarshal(msg.Body, &m)
						if err != nil {
							_ = msg.Reject(false)
							log.Println("items consumer got error: json.unmarshal msg.body: ", err)
						}

						switch m.Status {
						case model.ItemCreated:
							m.Status = model.ItemPending

							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								log.Println("items consumer got error: case: ItemCreated: UpdateItem: ", err)
								_ = msg.Reject(false)
								continue
							}
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error: case: ItemCreated: msg.Ack: ", err)
								continue
							}

						case model.ItemPending:
							m.Status = model.ItemComplete
							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								log.Println("items consumer got error: case: ItemPending: UpdateItem: ", err)
								_ = msg.Reject(false)
								continue
							}
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error: case: ItemPending: msg.Ack: ", err)
								continue
							}

						case model.ItemComplete:
							m.Status = model.ItemDeleted

							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								log.Println("items consumer got error: case: ItemComplete: UpdateItem: ", err)
								_ = msg.Reject(false)
								continue
							}
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error: case: ItemComplete: msg.Ack: ", err)
								continue
							}
							continue

						case model.ItemDeleted:
							// for next implemented usage
							err := msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error: case: ItemDeleted: msg.Ack: ", err)
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

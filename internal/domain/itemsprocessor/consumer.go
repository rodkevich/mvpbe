package itemsprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"

	"github.com/rodkevich/mvpbe/internal/domain/itemsprocessor/model"
)

func runExampleItemsConsumer(ctx context.Context, itemsUsage ItemsSampleProcessUsage, itemsCh <-chan amqp.Delivery) {
	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i <= exAMQPConcurrencyItems; i++ {
		eg.Go(func(ctx context.Context, itemsCh <-chan amqp.Delivery, workerID int) func() error {
			log.Printf("starting consumer id: %d, for items queue: %s", workerID, exQueueNameItems)

			return func() error {
				for {
					select {
					case <-ctx.Done():
						log.Printf("items consumer ctx done: %v", ctx.Err())
						return ctx.Err()

					case msg, ok := <-itemsCh:
						if !ok {
							log.Printf("error: items channel closed for queue: %s", exQueueNameItems)
							return errors.New("items channel closed")
						}

						log.Printf("Items consumer: id: %d, data: %s, headers: %+v", workerID, string(msg.Body), msg.Headers)

						m := model.SampleItem{}
						err := json.Unmarshal(msg.Body, &m)
						if err != nil {
							_ = msg.Reject(false)
							log.Println("items consumer got error: json.unmarshal msg.body: ", err)
						}

						// todo remove if no need to simulate
						fakeJobTime := 3 * time.Second

						switch m.Status {
						case model.ItemCreated:
							time.Sleep(fakeJobTime)

							m.Status = model.ItemPending
							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								log.Println("items consumer got error: case: ItemCreated: UpdateItem: ", err)
								_ = msg.Nack(false, true)
							}
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error: case: ItemCreated: msg.Ack: ", err)
							}

						case model.ItemPending:
							time.Sleep(fakeJobTime)

							m.Status = model.ItemComplete
							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								log.Println("items consumer got error: case: ItemPending: UpdateItem: ", err)
								_ = msg.Nack(false, true)
							}
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error: case: ItemPending: msg.Ack: ", err)
							}

						case model.ItemComplete:
							time.Sleep(fakeJobTime)

							m.Status = model.ItemDeleted
							err = itemsUsage.UpdateItem(ctx, &m)
							if err != nil {
								log.Println("items consumer got error: case: ItemComplete: UpdateItem: ", err)
								_ = msg.Nack(false, true)
							}
							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer got error: case: ItemComplete: msg.Ack: ", err)
							}

						case model.ItemDeleted:
							// shouldn't appear here anymore // todo remove after 'deleted' worker is done and tested
							err := msg.Nack(false, true)
							if err != nil {
								log.Println("items consumer got error: case: deleted item: msg.Nack: ", err)
							}
						default:
							// return to que
							_ = msg.Nack(false, true)
						}
					}
				}
			}
		}(ctx, itemsCh, i))
	}
	_ = eg.Wait()
}

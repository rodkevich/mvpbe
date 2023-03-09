package itemsprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"golang.org/x/sync/errgroup"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rodkevich/mvpbe/internal/domain/itemsprocessor/model"
)

func runExampleItemsConsumer(ctx context.Context, itemsUsage Processor, itemsCh <-chan amqp.Delivery) {
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

						log.Printf("Consumer id: %d, data: %s, headers: %+v", workerID, string(msg.Body), msg.Headers)

						task := &model.SomeProcessingTask{}

						err := json.Unmarshal(msg.Body, &task.SampleItem)
						if err != nil {
							_ = msg.Reject(false)
							log.Println("items consumer error: json.unmarshal msg.body: ", err)
						}

						if headerValue, ok := msg.Headers["example-item-trace-id"].(string); ok {
							task.TraceID = headerValue

							err = SaveState(task)
							if err != nil {
								log.Println("items consumer error: SaveState: ", err)
							}
						}

						// todo remove if no need to simulate
						fakeJobTime := 2 * time.Second

						switch task.Status {
						case model.ItemCreated:
							time.Sleep(8 * time.Second)

							task.Status = model.ItemPending
							err = itemsUsage.UpdateItem(ctx, task)
							if err != nil {
								log.Println("items consumer error: case: ItemCreated: UpdateItem: ", err)
								_ = msg.Nack(false, true)
							}

							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer error: case: ItemCreated: msg.Ack: ", err)
							}

						case model.ItemPending:
							time.Sleep(fakeJobTime)

							task.Status = model.ItemComplete
							err = itemsUsage.UpdateItem(ctx, task)
							if err != nil {
								log.Println("items consumer error: case: ItemPending: UpdateItem: ", err)
								_ = msg.Nack(false, true)
							}

							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer error: case: ItemPending: msg.Ack: ", err)
							}

						case model.ItemComplete:
							time.Sleep(fakeJobTime)

							task.Status = model.ItemDeleted
							err = itemsUsage.UpdateItem(ctx, task)
							if err != nil {
								log.Println("items consumer error: case: ItemComplete: UpdateItem: ", err)
								_ = msg.Nack(false, true)
							}

							err = msg.Ack(false)
							if err != nil {
								log.Println("items consumer error: case: ItemComplete: msg.Ack: ", err)
							}

							log.Printf("Total items having saved states: %d Item id [%s] saved states: %d", StatesLength(), task.TraceID, StatesLengthByID(task.TraceID))

						case model.ItemDeleted:
							// shouldn't appear here anymore // todo remove after 'deleted' worker is done and tested
							err := msg.Nack(false, true)
							if err != nil {
								log.Println("items consumer error: case: deleted item: msg.Nack: ", err)
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

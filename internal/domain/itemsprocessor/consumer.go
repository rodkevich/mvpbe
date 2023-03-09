package itemsprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"golang.org/x/sync/errgroup"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rodkevich/mvpbe/internal/domain/itemsprocessor/model"
)

func runExampleItemsConsumer(ctx context.Context, itemsUsage Processor, itemsCh <-chan amqp.Delivery) {
	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i <= exAMQPConcurrency; i++ {
		eg.Go(func(ctx context.Context, itemsCh <-chan amqp.Delivery, workerID int) func() error {
			log.Printf("starting consumer id: %d, for items queue: %s", workerID, exQueueNameProcess)

			return func() error {
				for {
					select {
					case <-ctx.Done():
						log.Printf("items consumer ctx done: %v", ctx.Err())
						return ctx.Err()

					case msg, ok := <-itemsCh:
						if !ok {
							log.Printf("error: items channel closed for queue: %s", exQueueNameProcess)
							return errors.New("items channel closed")
						}

						// Log incoming message and consumer attributes,
						// so as message body and incoming headers
						log.Printf("Consumer id: %d, data: %s, headers: %+v", workerID, string(msg.Body), msg.Headers)

						// Create new sample processing task for worker.
						// Unmarshal body into it. If body can't be processed - reject this message
						// to free que and prevent it's infinite processing in future
						task := &model.SomeProcessingTask{}
						err := json.Unmarshal(msg.Body, &task.SampleItem)
						if err != nil {
							_ = msg.Reject(false)
							log.Println("items consumer error: json.unmarshal msg.body: ", err)
						}

						// Try to get trace id from message headers.
						// If no valid identifier presented - skip saving
						// of items states to dispatcher.
						if headerValue, ok := msg.Headers["example-item-trace-id"].(string); ok {
							task.TraceID = headerValue

							// If valid trace id comes, save task state.
							// Saved states can possibly be used for undo/redo operations.
							err = SaveState(task)
							if err != nil {
								log.Println("items consumer error: SaveState: ", err)
							}
						}

						// For demo purpose imitate some job around tasks
						// according their status field.
						switch task.Status {
						case model.ItemCreated:
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

							log.Printf("Total items having saved states: %d Item id [%s] saved states: %d",
								StatesLength(), task.TraceID, StatesLengthByID(task.TraceID))

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

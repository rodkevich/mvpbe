package itemsprocessor

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rodkevich/mvpbe/internal/domain/itemsprocessor/model"
)

func init() {
	go dispatcher.backgroundExpire()
}

const initialSize = 16

// let it be singleton
var dispatcher = &stateDispatcher{
	mu:     sync.RWMutex{},
	tasks:  make(map[string][]*model.SampleItem, initialSize),
	stop:   make(chan bool),
	ticker: time.NewTicker(20 * time.Second),
}

type stateDispatcher struct {
	mu     sync.RWMutex
	tasks  map[string][]*model.SampleItem
	stop   chan bool
	ticker *time.Ticker
}

// StopDispatcher shut down the background cleanup
func StopDispatcher() {
	log.Println("Closing items states dispatcher.")
	dispatcher.ticker.Stop()
	dispatcher.stop <- true
}

func (c *stateDispatcher) backgroundExpire() {
	for {
		select {
		case <-c.stop:
			close(c.stop)
			return
		case t := <-c.ticker.C:
			log.Println("[background] Running  expire check for items")

			c.mark(t.UnixNano())
		}
	}
}

func (c *stateDispatcher) mark(_ int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for k, v := range c.tasks {
		name := k
		for index, i := range v {
			item := i
			if item.Status == model.ItemComplete && item.Expired() {
				log.Printf(
					"[background] Mark item states for deletition: name [%s], state [%s], time [%s] \n",
					name, item.Status, item.FinishTime)

				go c.purgeExpired(name, index, item.FinishTime.Unix())
			}
		}
	}
}

func (c *stateDispatcher) purgeExpired(name string, p int, expectedExpiryTime int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if items, ok := c.tasks[name]; ok && items[p].FinishTime.Unix() == expectedExpiryTime {
		log.Println("[background] Purge expired items", name, p, expectedExpiryTime)

		delete(c.tasks, name)
	}
}

func set(bs *stateDispatcher, key string, t *model.SampleItem) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.tasks[key] = append(bs.tasks[key], t)
	return nil
}

func get(bs *stateDispatcher, key string, index int) (*model.SampleItem, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if v, ok := bs.tasks[key]; ok {
		return v[index], true
	}

	return nil, false
}

func SaveState(t *model.SomeProcessingTask) error {
	if t.TraceID == "" || dispatcher.tasks == nil {
		return fmt.Errorf("unable to proceed: key: [%v], item: [%#v]", t.TraceID, t)
	}
	return set(dispatcher, t.TraceID, &t.SampleItem)
}

func GetState(key string, index int) (*model.SampleItem, bool) {
	if key == "" || dispatcher.tasks == nil {
		return nil, false
	}

	return get(dispatcher, key, index)
}

// ClearDispatcherStates remove all items
func ClearDispatcherStates() {
	dispatcher.mu.Lock()
	defer dispatcher.mu.Unlock()

	dispatcher.tasks = make(map[string][]*model.SampleItem, initialSize)
}

func StatesLength() int {
	dispatcher.mu.RLock()
	defer dispatcher.mu.RUnlock()
	return len(dispatcher.tasks)
}

func StatesLengthByID(key string) int {
	dispatcher.mu.RLock()
	defer dispatcher.mu.RUnlock()
	return len(dispatcher.tasks[key])
}

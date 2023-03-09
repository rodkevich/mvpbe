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
	ticker: time.NewTicker(5 * time.Second),
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

// StopDispatcher shut down the background cleanup
func StopDispatcher() {
	log.Println("Closing items states dispatcher.")
	dispatcher.ticker.Stop()
	dispatcher.stop <- true
}

func SaveState(m *model.SomeProcessingTask) error {
	if m.TraceID == "" || dispatcher.tasks == nil {
		return fmt.Errorf("unable to proceed: key: [%v], item: [%#v]", m.TraceID, m)
	}
	return set(dispatcher, m.TraceID, &m.SampleItem)
}

func GetState(key string, index int) (*model.SampleItem, bool) {
	if key == "" || dispatcher.tasks == nil {
		return nil, false
	}
	return get(dispatcher, key, index)
}

type stateDispatcher struct {
	mu     sync.RWMutex
	tasks  map[string][]*model.SampleItem
	stop   chan bool
	ticker *time.Ticker
}

func (sd *stateDispatcher) backgroundExpire() {
	for {
		select {
		case <-sd.stop:
			close(sd.stop)
			return
		case t := <-sd.ticker.C:
			log.Println("[background] Running expire check for items")
			sd.mark(t.UnixNano())
		}
	}
}

func (sd *stateDispatcher) mark(_ int64) {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	for k, v := range sd.tasks {
		key := k
		for index, i := range v {
			item := i
			if item.Status == model.ItemComplete && item.Expired() {
				log.Printf(
					"[background] Mark item states for deletition: name [%s], state [%s], time [%s] \n",
					key, item.Status, item.FinishTime)

				go sd.purgeExpired(key, index, item.FinishTime.Unix())
			}
		}
	}
}

func (sd *stateDispatcher) purgeExpired(key string, index int, time int64) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if items, ok := sd.tasks[key]; ok && items[index].FinishTime.Unix() == time {
		log.Println("[background] Purge expired items", key, index, time)

		delete(sd.tasks, key)
	}
}

func set(sd *stateDispatcher, key string, t *model.SampleItem) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	sd.tasks[key] = append(sd.tasks[key], t)
	return nil
}

func get(sd *stateDispatcher, key string, index int) (*model.SampleItem, bool) {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	if v, ok := sd.tasks[key]; ok {
		return v[index], true
	}
	return nil, false
}

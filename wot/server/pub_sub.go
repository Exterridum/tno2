package server

import (
	"sync"

	"github.com/conas/tno2/util/async"
)

type Subscribers struct {
	rwmut        *sync.RWMutex
	subscription map[string]*async.FanOut
}

func NewSubscribers() *Subscribers {
	return &Subscribers{
		rwmut:        &sync.RWMutex{},
		subscription: make(map[string]*async.FanOut),
	}
}

func (wss *Subscribers) CreateSubscription(subscriptionID string, clients *async.FanOut) {
	wss.rwmut.Lock()
	defer wss.rwmut.Unlock()

	wss.subscription[subscriptionID] = clients
}

func (wss *Subscribers) CancelSubscription(subscriptionID string) {
	wss.rwmut.Lock()
	defer wss.rwmut.Unlock()

	wss.subscription[subscriptionID].RemoveAllSubscribes()
	delete(wss.subscription, subscriptionID)
}

func (wss *Subscribers) AddClient(subscriptionID string, client chan<- interface{}) int {
	wss.rwmut.RLock()
	defer wss.rwmut.RUnlock()

	return wss.subscription[subscriptionID].AddSubscriber(client)
}

func (wss *Subscribers) RemoveClient(subscriptionID string, clientID int) {
	wss.rwmut.RLock()
	defer wss.rwmut.RUnlock()

	wss.subscription[subscriptionID].RemoveSubscriber(clientID)
}

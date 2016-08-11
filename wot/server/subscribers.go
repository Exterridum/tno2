package server

import (
	"sync"

	"github.com/conas/tno2/util/async"
)

// Subscribers struct allows to share one subscription between multiple clients
// Common scenario is when user creates subscription for event or fires action, link to subscription is returned
// and then multiple clients can share this subscription link
// Each entry in subscription map of Subscribers struct, corresponds to one real subscription. Map entry
// then contains all connected clients
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

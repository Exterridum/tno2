package async

import "sync"

//TODO make performance and memmory tests to decide between go routines and mutexes

type FanOut struct {
	out   map[string]chan<- interface{}
	mutex *sync.Mutex
}

func NewFanOut() *FanOut {
	return &FanOut{
		out:   make(map[string]chan<- interface{}),
		mutex: &sync.Mutex{},
	}
}

func (fo *FanOut) AddSubscriber(id string, out chan<- interface{}) {
	fo.mutex.Lock()
	fo.out[id] = out
	fo.mutex.Unlock()
}

func (fo *FanOut) RemoveSubscriber(id string) {
	fo.mutex.Lock()
	delete(fo.out, id)
	fo.mutex.Unlock()
}

func (fo *FanOut) Publish(event interface{}) {
	go func() {
		fo.mutex.Lock()
		for _, out := range fo.out {
			out <- event
		}
		fo.mutex.Unlock()
	}()
}

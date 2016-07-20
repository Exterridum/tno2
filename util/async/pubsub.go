package async

import "sync"

//TODO make performance and memmory tests to decide between go routines and mutexes

type FanOut struct {
	out   map[int]chan<- interface{}
	mutex *sync.RWMutex
}

func NewFanOut() *FanOut {
	return &FanOut{
		out:   make(map[int]chan<- interface{}),
		mutex: &sync.RWMutex{},
	}
}

func (fo *FanOut) AddSubscriber(out chan<- interface{}) int {
	fo.mutex.Lock()
	defer fo.mutex.Unlock()

	size := len(fo.out)
	fo.out[size] = out

	return size
}

func (fo *FanOut) RemoveSubscriber(id int) {
	fo.mutex.Lock()
	defer fo.mutex.Unlock()

	close(fo.out[id])
	delete(fo.out, id)
}

func (fo *FanOut) RemoveAllSubscribes() {
	fo.mutex.Lock()
	defer fo.mutex.Unlock()

	fo.out = make(map[int]chan<- interface{})
}

func (fo *FanOut) Publish(event interface{}) {
	go func() {
		fo.mutex.RLock()
		defer fo.mutex.RUnlock()

		for _, out := range fo.out {
			out <- event
		}
	}()
}

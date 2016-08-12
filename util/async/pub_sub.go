package async

import "sync"

type FanOut struct {
	out     map[int]chan<- interface{}
	l       *sync.RWMutex
	counter int
	pool    []int
}

func NewFanOut() *FanOut {
	return &FanOut{
		out:  make(map[int]chan<- interface{}),
		l:    &sync.RWMutex{},
		pool: make([]int, 0),
	}
}

func (fo *FanOut) AddSubscriber(out chan<- interface{}) int {
	fo.l.Lock()

	id := fo.nextID()
	fo.out[id] = out

	fo.l.Unlock()
	return id
}

func (fo *FanOut) nextID() int {
	//use id form pool if any available
	if len(fo.pool) > 0 {
		id := fo.pool[0]
		fo.pool = deleteElement(fo.pool, 0)
		return id
	}

	for {
		_, exists := fo.out[fo.counter]
		if !exists {
			break
		}
		fo.counter++
	}

	return fo.counter
}

func (fo *FanOut) RemoveSubscriber(id int) {
	fo.l.Lock()
	//TODO: investigate if close on channel can cause panic. If not move Unlock from defered
	//to end of the method for performance reasons
	defer fo.l.Unlock()

	if _, ok := fo.out[id]; ok {
		delete(fo.out, id)
		fo.pool = append(fo.pool, id)
	}
}

func (fo *FanOut) RemoveAllSubscribes() {
	fo.l.Lock()
	//TODO: investigate if close on channel can cause panic. If not move Unlock from defered
	//to end of the method for performance reasons
	defer fo.l.Unlock()

	fo.out = make(map[int]chan<- interface{})
	fo.pool = make([]int, 0)
	fo.counter = 0
}

func (fo *FanOut) Publish(event interface{}) {
	go func() {
		fo.l.RLock()
		outCopy := mapClone(fo.out)
		fo.l.RUnlock()

		//non blocking message publish
		for {
			if len(outCopy) == 0 {
				break
			}

			for k, out := range outCopy {
				//TODO: should be publishing time limited so we break cycle even of not
				//all subscribers were not sent messages to?
				select {
				case out <- event:
					//TODO: What is the performance of remove element? If expensive implement
					//another solution to tag published subscribers
					delete(outCopy, k)
				}
			}
		}
	}()
}

func mapClone(src map[int]chan<- interface{}) map[int]chan<- interface{} {
	newMap := make(map[int]chan<- interface{})

	for k, v := range src {
		newMap[k] = v
	}

	return newMap
}

func deleteElement(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

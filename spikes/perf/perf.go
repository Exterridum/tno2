package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	chanTest()
	bufferedChanTest()
	goRoutineTest()
	callbackTest()
	callbackMutexTest()
	callbackMutexDeferTest()
	mapFuncTest()
	atomicUpdateTest()
}

var sampleSize = 10000000

func chanTest() {
	c := make(chan int)

	go func() {
		for {
			<-c
		}
	}()

	time.Sleep(time.Second * 1)

	defer timeTrack(time.Now(), "channels", sampleSize)
	for i := 0; i < sampleSize; i++ {
		c <- i
	}
}

func bufferedChanTest() {
	c := make(chan int, 10)

	go func() {
		for {
			<-c
		}
	}()

	time.Sleep(time.Second * 1)

	defer timeTrack(time.Now(), "buffered channels", sampleSize)
	for i := 0; i < sampleSize; i++ {
		c <- i
	}
}

func goRoutineTest() {
	defer timeTrack(time.Now(), "go routine", sampleSize)

	for i := 0; i < sampleSize; i++ {
		go func() {}()
	}
}

func callbackTest() {
	defer timeTrack(time.Now(), "callback simple", sampleSize)

	for i := 0; i < sampleSize; i++ {
		callbackSimple()
	}
}

func callbackSimple() {
}

func callbackMutexTest() {
	defer timeTrack(time.Now(), "callback mutex", sampleSize)

	for i := 0; i < sampleSize; i++ {
		callbackMutex()
	}
}

var m = &sync.Mutex{}

func callbackMutex() {
	m.Lock()
	m.Unlock()
}

func callbackMutexDeferTest() {
	defer timeTrack(time.Now(), "callback mutex defer", sampleSize)

	for i := 0; i < sampleSize; i++ {
		callbackDeferMutex()
	}
}

func callbackDeferMutex() {
	m.Lock()
	defer m.Unlock()
}

func mapFuncTest() {
	m := make(map[int]func(int) int)

	m[1] = fn1
	m[2] = fn2
	m[3] = fn3
	m[4] = fn4
	m[5] = fn5

	defer timeTrack(time.Now(), "map function test", sampleSize)
	for i := 0; i < sampleSize/5; i++ {
		m[1](i)
		m[2](i)
		m[3](i)
		m[4](i)
		m[5](i)
	}

}

func fn1(v int) int {
	return 1
}

func fn2(v int) int {
	return 2
}
func fn3(v int) int {
	return 3
}
func fn4(v int) int {
	return 4
}
func fn5(v int) int {
	return 5
}

func atomicUpdateTest() {
	v := atomic.Value{}

	defer timeTrack(time.Now(), "atomic update test", sampleSize)
	for i := 0; i < sampleSize; i++ {
		v.Store(i)
	}
}

func timeTrack(start time.Time, name string, op int) {
	elapsed := time.Since(start)
	log.Printf("%20s time: %15s, ops/s: %v", name, elapsed, float64(op)/elapsed.Seconds())
}

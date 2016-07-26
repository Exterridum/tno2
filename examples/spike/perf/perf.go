package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	chanTest()
	bufferedChanTest()
	goRoutineTest()
	callbackTest()
	callbackMutexTest()
	callbackMutexDeferTest()
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

	defer timeTrack(time.Now(), "channels")
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

	defer timeTrack(time.Now(), "buffered channels")
	for i := 0; i < sampleSize; i++ {
		c <- i
	}
}

func goRoutineTest() {
	defer timeTrack(time.Now(), "go routine")

	for i := 0; i < sampleSize; i++ {
		go func() {}()
	}
}

func callbackTest() {
	defer timeTrack(time.Now(), "callback simple")

	for i := 0; i < sampleSize; i++ {
		callbackSimple()
	}
}

func callbackSimple() {
}

func callbackMutexTest() {
	defer timeTrack(time.Now(), "callback mutex")

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
	defer timeTrack(time.Now(), "callback mutex defer")

	for i := 0; i < sampleSize; i++ {
		callbackDeferMutex()
	}
}

func callbackDeferMutex() {
	m.Lock()
	defer m.Unlock()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

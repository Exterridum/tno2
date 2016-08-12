package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type args struct {
	Foo, Bar string
}

var sampleSize = 100000

func main() {
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	defer timeTrack(time.Now(), "callback mutex defer", sampleSize)

	for i := 0; i < sampleSize; i++ {
		// e := client.Call("Compose.Details", &args{"Foo!", "Bar!"}, &reply)
		client.Call("Compose.Details", &args{"Foo!", "Bar!"}, &reply)
	}
	// if e != nil {
	// 	log.Fatalf("Something went wrong: %v", e.Error())
	// }

	fmt.Printf("The 'reply' pointer value has been changed to: %s", reply)
}

func timeTrack(start time.Time, name string, op int) {
	elapsed := time.Since(start)
	log.Printf("%20s time: %15s, ops/s: %v", name, elapsed, float64(op)/elapsed.Seconds())
}

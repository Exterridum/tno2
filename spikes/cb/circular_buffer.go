package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	buff := NewBuffer(16)

	c1 := buff.Client()
	c2 := buff.Client()

	wg := sync.WaitGroup{}
	wg.Add(2)

	dataSize := 1000000

	defer timeTrack(time.Now(), "atomic update test", dataSize)

	go func() {
		for i := 0; i < dataSize; i++ {
			v := &Val{
				empty: false,
				data:  i,
			}
			c1.Write(v)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < dataSize; i++ {
			c2.Read()
		}
		wg.Done()
	}()

	wg.Wait()
}

func timeTrack(start time.Time, name string, op int) {
	elapsed := time.Since(start)
	log.Printf("%20s time: %15s, ops/s: %v", name, elapsed, float64(op)/elapsed.Seconds())
}

type Buffer struct {
	size int
	buff []*atomic.Value
}

func NewBuffer(size int) Buffer {
	b := make([]*atomic.Value, size)

	for i, _ := range b {
		b[i] = &atomic.Value{}
		b[i].Store(emptySlot)
	}

	return Buffer{
		size: size,
		buff: b,
	}
}

func (b *Buffer) Get(i int) interface{} {
	return b.buff[i].Load()
}

func (b *Buffer) Set(i int, val interface{}) {
	b.buff[i].Store(val)
}

type Val struct {
	empty bool
	data  interface{}
}

var emptySlot = &Val{
	empty: true,
	data:  nil,
}

type Client struct {
	buf  *Buffer
	size int
	pos  int
}

func (b *Buffer) Client() *Client {
	return &Client{
		buf:  b,
		size: b.size,
		pos:  b.size - 1,
	}
}

func (c *Client) Read() interface{} {
	for {
		v := c.buf.Get(c.pos)
		if v.(*Val).empty == true {
			continue
		}
		c.buf.Set(c.pos, emptySlot)
		c.pos = (c.pos + 1) % c.size
		return v
	}
}

func (c *Client) Write(nv interface{}) {
	for {
		v := c.buf.Get(c.pos)
		if v.(*Val).empty == false {
			continue
		}
		c.buf.Set(c.pos, nv)
		c.pos = (c.pos + 1) % c.size
		return
	}
}

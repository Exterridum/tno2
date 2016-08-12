package async

import (
	"testing"

	"github.com/conas/tno2/util/str"
)

func TestCaseFanOutIdRecycle(t *testing.T) {
	fo := NewFanOut()
	assertFanOut("FanOutIdRecycle.0", t, fo, 0, 0, 0)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.1", t, fo, 0, 0, 1)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.2", t, fo, 1, 0, 2)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.3", t, fo, 2, 0, 3)

	fo.RemoveSubscriber(1)
	assertFanOut("FanOutIdRecycle.4", t, fo, 2, 1, 2)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.5", t, fo, 2, 0, 3)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.6", t, fo, 3, 0, 4)

	fo.RemoveSubscriber(1)
	assertFanOut("FanOutIdRecycle.7", t, fo, 3, 1, 3)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.8", t, fo, 3, 0, 4)

	fo.RemoveSubscriber(2)
	assertFanOut("FanOutIdRecycle.9", t, fo, 3, 1, 3)

	fo.RemoveSubscriber(3)
	assertFanOut("FanOutIdRecycle.10", t, fo, 3, 2, 2)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.11", t, fo, 3, 1, 3)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.12", t, fo, 3, 0, 4)

	fo.AddSubscriber(make(chan interface{}))
	assertFanOut("FanOutIdRecycle.13", t, fo, 4, 0, 5)

	fo.RemoveAllSubscribes()
	assertFanOut("FanOutIdRecycle.14", t, fo, 0, 0, 0)
}

func assertFanOut(msg string, t *testing.T, fo *FanOut, id, poolLen, outLen int) {
	Equals(str.Concat(msg, " fo.id"), t, id, fo.counter)
	Equals(str.Concat(msg, " len(fo.pool)"), t, poolLen, len(fo.pool))
	Equals(str.Concat(msg, " len(fo.out)"), t, outLen, len(fo.out))
}

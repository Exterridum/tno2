package async

import (
	"testing"
	"time"
)

func TestCase1(t *testing.T) {
	Run(func() interface{} {
		return 3 + 4
	}).Then(func(val interface{}) interface{} {
		return 2 * val.(int)
	}).Then(func(val interface{}) interface{} {
		Equals(t, 14, val.(int))
		return nil
	})

	time.Sleep(time.Second * 1)
}

func Equals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Log("\nExpected:", expected, "\nReturned:", actual)
		t.Fail()
	}
}

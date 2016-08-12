package async

import "testing"

func TestCase1(t *testing.T) {
	r := Run(func() interface{} {
		return 3 + 4
	}).Then(func(val interface{}) interface{} {
		return 2 * val.(int)
	}).Then(func(val interface{}) interface{} {
		return val
	}).Get()

	Equals("TestCase1", t, 14, r.(int))
}

func Equals(assetName string, t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Log("\nTest: ", assetName, "\nExpected:", expected, "\nReturned:", actual)
		t.Fail()
	}
}

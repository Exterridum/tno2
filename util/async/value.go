package async

// ----- USING CHANNEL

type Value struct {
	out chan interface{}
}

func newValue() *Value {
	return &Value{
		out: make(chan interface{}, 1),
	}
}

// Get blocks until value is available
func (v *Value) Get() interface{} {
	return <-v.out
}

func (v *Value) set(data interface{}) {
	v.out <- data
}

// ----- USING WAIT GROUP
// type Value struct {
// 	wg   *sync.WaitGroup
// 	data interface{}
// }

// func newValue() *Value {
// 	wg := &sync.WaitGroup{}
// 	wg.Add(1)

// 	return &Value{
// 		wg: wg,
// 	}
// }

// func (v *Value) Get() interface{} {
// 	v.wg.Wait()
// 	return v.data
// }

// func (v *Value) set(data interface{}) {
// 	v.data = data
// 	v.wg.Done()
// }

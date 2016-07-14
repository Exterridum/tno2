package driver

type Driver struct {
	in <-chan interface{}
}

func ReplyDriver() *Driver {
	return &Driver{}
}

func (d *Driver) SetInputChannel(ich <-chan interface{}) {
	d.in = ich
}

func (d *Driver) Start() {
	go func() {
		for msg := range d.in {
			handle(msg)
		}
	}()
}

func handle(msg interface{}) {

}

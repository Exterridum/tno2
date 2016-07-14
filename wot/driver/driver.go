package driver

type Driver interface {
	SetInputChannel(in <-chan interface{})
	Start()
}

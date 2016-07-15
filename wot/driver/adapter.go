package driver

import "github.com/conas/tno2/concurrent"

type Adapter struct {
	driver Driver
}

func NewAdapter(driver Driver) *Adapter {
	return &Adapter{
		driver: driver,
	}
}

func (a *Adapter) Send(msg Message) *concurent.Promise {
	p := concurent.NewPromise()

	switch msg.GetMessageType() {
	case IARQ:
		go func() {
			p.Channel() <- a.driver.InvokeAction(msg.(*InvokeActionRQ))
		}()

		return p
	case GPRQ:
		go func() {
			p.Channel() <- a.driver.GetProperty(msg.(*GetPropertyRQ))
		}()

		return p
	case SPRQ:
		go func() {
			a.driver.SetProperty(msg.(*SetPropertyRQ))
			p.Channel() <- nil
		}()

		return p
	default:
		return p
	}
}

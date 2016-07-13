package driver

type InvokeActionRQ struct {
	Response  chan interface{}
	Name      string
	Parameter interface{}
}

type GetPropertyRQ struct {
	Name string
}

type GetPropertyRS struct {
	Value interface{}
}

type SetPropertyRQ struct {
	Name  string
	Value interface{}
}

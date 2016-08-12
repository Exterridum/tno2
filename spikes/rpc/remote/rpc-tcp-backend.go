package main

// TCPArgs is structured around the client's provided parameters
// The struct's fields need to be exported too!
type TCPArgs struct {
	Foo string
	Bar string
}

// Compose is our RPC functions return type
type Compose string

// Details is our exposed RPC function
func (c *Compose) Details(args *TCPArgs, reply *string) error {
	*c = "some value"
	*reply = "Blah!"
	return nil
}

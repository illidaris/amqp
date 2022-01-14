package amqp

type Element struct {
	Name       string
	Durable    bool
	AutoDelete bool
	NoWait     bool
	Arguments  map[string]interface{}
}

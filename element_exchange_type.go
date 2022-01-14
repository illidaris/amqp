package amqp

type ExchangeType int32

const (
	Nil ExchangeType = iota
	Direct
	Topic
	FanOut
	Headers
)

func (i ExchangeType) String() string {
	switch i {
	case Direct:
		return "direct"
	case Topic:
		return "topic"
	case FanOut:
		return "fanout"
	case Headers:
		return "headers"
	default:
		return "nil"
	}
}

package amqp

import (
	amqpMeta "github.com/streadway/amqp"
)

type Exchange struct {
	Element
	Kind     ExchangeType
	Internal bool
}

func (e *Exchange) Declare(channel *amqpMeta.Channel) error {
	return channel.ExchangeDeclare(e.Name, e.Kind.String(), e.Durable, e.AutoDelete, e.Internal, e.NoWait, e.Arguments)
}

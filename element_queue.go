package amqp

import (
	amqpMeta "github.com/streadway/amqp"
)

type Queue struct {
	Element
	Exclusive bool
}

func (e *Queue) Declare(channel *amqpMeta.Channel) (amqpMeta.Queue, error) {
	return channel.QueueDeclare(e.Name, e.Durable, e.AutoDelete, e.Exclusive, e.NoWait, e.Arguments)
}

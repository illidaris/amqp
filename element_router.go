package amqp

import (
	amqpMeta "github.com/streadway/amqp"
)

type Router struct {
	Ex        *Exchange
	Q         *Queue
	Name      string
	NoWait    bool
	Arguments map[string]interface{}
}

func (e *Router) Bind(channel *amqpMeta.Channel) error {
	return channel.QueueBind(e.Q.Name, e.Name, e.Ex.Name, e.NoWait, e.Arguments)
}

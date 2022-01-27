package amqp

import (
	"context"
	amqpMeta "github.com/streadway/amqp"
)

type IConsumer interface {
	ICaller
}

type Consumer struct {
	Name             string
	Ctx              context.Context
	QueueName        string
	NoLocal          bool
	AutoAck          bool
	Exclusive        bool
	NoWait           bool
	Arguments        map[string]interface{}
	CloseHandlers    []CloseHandler
	DeliveryHandlers []DeliveryHandler
	Channel          *amqpMeta.Channel

	connCloseCh    <-chan *amqpMeta.Error
	channelCloseCh <-chan *amqpMeta.Error
}

func (e *Consumer) onClose(err *amqpMeta.Error) {
	if e.CloseHandlers != nil {
		for _, h := range e.CloseHandlers {
			h(err)
		}
	}
}

func (e *Consumer) onDelivery(delivery amqpMeta.Delivery) {
	if e.DeliveryHandlers != nil {
		for _, h := range e.DeliveryHandlers {
			h(delivery)
		}
	}
}

func (e *Consumer) Register(m *AMQPManager) error {
	// TODO: reconnect design
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err := e.Channel.Qos(1, 0, true)
	if err != nil {
		return err
	}
	deliveryCh, err := e.Channel.Consume(e.QueueName, e.Name, e.AutoAck, e.Exclusive, e.NoLocal, e.NoWait, e.Arguments)
	if err != nil {
		return err
	}
	conn, err := m.GetConnect()
	if err != nil {
		return err
	}
	e.connCloseCh = conn.NotifyClose(make(chan *amqpMeta.Error))
	e.channelCloseCh = e.Channel.NotifyClose(make(chan *amqpMeta.Error))

	m.consumers = append(m.consumers, e)
	go func() {
		defer e.Channel.Close()
		for {
			select {
			case closeErr := <-e.connCloseCh:
				if err != nil {
					e.onClose(closeErr)
				}
				return
			case closeErr := <-e.channelCloseCh:
				if err != nil {
					e.onClose(closeErr)
				}
				return
			case delivery := <-deliveryCh:
				if delivery.Body != nil {
					e.onDelivery(delivery)
				}
			}
		}
	}()

	return err
}

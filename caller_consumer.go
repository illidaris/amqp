package amqp

import (
	"context"
	amqpMeta "github.com/streadway/amqp"
)

type IConsumer interface {
	ICaller
	AddCloseHandler(h CloseHandler)
	AddDeliveryHandler(h DeliveryHandler)
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

	m              *AMQPManager
	connCloseCh    <-chan *amqpMeta.Error
	channelCloseCh <-chan *amqpMeta.Error
}

func NewConsumer(ctx context.Context, name, queue string) IConsumer {
	return &Consumer{
		Name:             name,
		Ctx:              ctx,
		QueueName:        queue,
		NoLocal:          false,
		AutoAck:          false,
		Exclusive:        false,
		NoWait:           false,
		CloseHandlers:    make([]CloseHandler, 0),
		DeliveryHandlers: make([]DeliveryHandler, 0),
	}
}

func (e *Consumer) AddCloseHandler(h CloseHandler) {
	e.CloseHandlers = append(e.CloseHandlers, h)
}

func (e *Consumer) AddDeliveryHandler(h DeliveryHandler) {
	e.DeliveryHandlers = append(e.DeliveryHandlers, h)
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
	conn, err := m.GetConnect()
	if err != nil {
		return err
	}
	e.connCloseCh = conn.NotifyClose(make(chan *amqpMeta.Error))
	e.m = m
	m.consumers[e.Name] = e
	ch, err := m.NewChannel()
	if err != nil {
		return err
	}
	e.channelCloseCh = ch.NotifyClose(make(chan *amqpMeta.Error))
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = ch.Qos(1, 0, true)
	if err != nil {
		return err
	}
	deliveryCh, err := ch.Consume(e.QueueName, e.Name, e.AutoAck, e.Exclusive, e.NoLocal, e.NoWait, e.Arguments)
	if err != nil {
		return err
	}
	go func() {
		defer ch.Close()
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

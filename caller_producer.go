package amqp

import (
	"context"
	amqpMeta "github.com/streadway/amqp"
	"sync"
)

var _ = IProducer(&Producer{})

type IProducer interface {
	ICaller
	Publish(ctx context.Context, router string, encoder Encoder, messages ...interface{}) ([]error, error)
}

type Producer struct {
	Exchange      string
	lock          sync.RWMutex
	m             *AMQPManager
	closeHandlers []CloseHandler
	connCloseCh   <-chan *amqpMeta.Error
}

func NewProducer(exchange string) IProducer {
	p := &Producer{
		Exchange:      exchange,
		closeHandlers: make([]CloseHandler, 0),
	}
	return p
}

func (p *Producer) onClose(err *amqpMeta.Error) {
	if p.closeHandlers != nil {
		for _, h := range p.closeHandlers {
			h(err)
		}
	}
}

func (p *Producer) Register(m *AMQPManager) error {
	conn, err := m.GetConnect()
	if err != nil {
		return err
	}
	p.connCloseCh = conn.NotifyClose(make(chan *amqpMeta.Error))
	m.producers = append(m.producers, p)
	return nil
}

// Publish send message
func (p *Producer) Publish(ctx context.Context, router string, encoder Encoder, messages ...interface{}) ([]error, error) {
	// TODO: reconnect design
	// create new channel
	ch, err := p.m.NewChannel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()
	return send(ctx, ch, p.Exchange, router, encoder, messages...)
}

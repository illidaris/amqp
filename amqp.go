package amqp

import (
	"context"
	"fmt"
	amqpMeta "github.com/streadway/amqp"
	"sync"
)

// AMQPManager mq manager
type AMQPManager struct {
	lock sync.RWMutex
	TCPSection

	connection *amqpMeta.Connection
	channel    *amqpMeta.Channel
	producers  []IProducer
	consumers  []IConsumer
}

// URL build url string
func (m *AMQPManager) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		m.GetUser(),
		m.GetPwd(),
		m.GetHost(),
		m.GetPort(),
		m.GetPath(),
	)
}

// Connect open connection
func (m *AMQPManager) Connect() error {
	conn, err := amqpMeta.Dial(m.URL())
	if err != nil {
		return err
	}
	m.connection = conn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	m.channel = ch
	return nil
}

// DisConnect close connection
func (m *AMQPManager) DisConnect() error {
	defer m.connection.Close()
	return m.channel.Close()
}

// NewChannel get new channel in a living connect
func (m *AMQPManager) NewChannel() (*amqpMeta.Channel, error) {
	return m.connection.Channel()
}

// Declare declare some element, such as exchange/queue/router
func (m *AMQPManager) Declare(declareFunc ...DeclareFunc) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, f := range declareFunc {
		err := f(m.channel)
		if err != nil {
			return err
		}
	}
	return nil
}

// Register register consumer/producer in manager
func (m *AMQPManager) Register(c ICaller) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return c.Register(m)
}

// PublishOnce publish message in new connect
func (m *AMQPManager) PublishOnce(ctx context.Context, exchange, router string, encoder Encoder, messages ...interface{}) ([]error, error) {
	conn, err := amqpMeta.Dial(m.URL())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()
	return send(ctx, ch, exchange, router, encoder, messages...)
}

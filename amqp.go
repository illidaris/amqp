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
	producers  []IProducer
	consumers  []IConsumer
}

func NewManager(host, user, pwd, path string, port int32) *AMQPManager {
	manager := &AMQPManager{
		TCPSection: TCPSection{},
	}
	manager.SetHost(host)
	manager.SetPort(port)
	manager.SetUser(user)
	manager.SetPwd(pwd)
	manager.SetPath(path)
	return manager
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

// GetConnect get open connection
func (m *AMQPManager) GetConnect() (*amqpMeta.Connection, error) {
	if m.connection == nil || m.connection.IsClosed() {
		m.lock.Lock()
		defer m.lock.Unlock()
		conn, err := amqpMeta.Dial(m.URL())
		if err != nil {
			return nil, err
		}
		m.connection = conn
	}
	return m.connection, nil
}

// DisConnect close connection
func (m *AMQPManager) DisConnect() error {
	return m.connection.Close()
}

// NewChannel get new channel in a living connect
func (m *AMQPManager) NewChannel() (*amqpMeta.Channel, error) {
	conn, err := m.GetConnect()
	if err != nil {
		return nil, err
	}
	return conn.Channel()
}

// Declare declare some element, such as exchange/queue/router
func (m *AMQPManager) Declare(channel *amqpMeta.Channel, declareFunc ...DeclareFunc) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, f := range declareFunc {
		err := f(channel)
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
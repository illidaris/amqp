package amqp

import (
	"context"
	"fmt"
	amqpMeta "github.com/streadway/amqp"
	"sync"
	"time"
)

// AMQPManager mq manager
type AMQPManager struct {
	connLock     sync.RWMutex
	registerLock sync.RWMutex
	TCPSection
	connection   *amqpMeta.Connection
	producers    map[string]IProducer
	consumers    map[string]IConsumer
	closeCh      chan *amqpMeta.Error
	noAutoRelink chan struct{}
}

func NewManager(host, user, pwd, path string, port int32) *AMQPManager {
	manager := &AMQPManager{
		TCPSection:   TCPSection{},
		producers:    make(map[string]IProducer),
		consumers:    make(map[string]IConsumer),
		noAutoRelink: make(chan struct{}, 1),
		//closeCh: make(chan *amqpMeta.Error,1),
	}
	manager.SetHost(host)
	manager.SetPort(port)
	manager.SetUser(user)
	manager.SetPwd(pwd)
	manager.SetPath(path)
	return manager
}

func (m *AMQPManager) GetProducer(name string) IProducer {
	return m.producers[name]
}

func (m *AMQPManager) GetConsumer(name string) IConsumer {
	return m.consumers[name]
}

func (m *AMQPManager) Producers() map[string]IProducer {
	return m.producers
}

func (m *AMQPManager) Consumers() map[string]IConsumer {
	return m.consumers
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
	m.connLock.Lock()
	defer m.connLock.Unlock()
	if m.connection == nil || m.connection.IsClosed() {
		ch := make(chan *amqpMeta.Error, 1)
		m.closeCh = ch
		conn, err := amqpMeta.Dial(m.URL())
		if err != nil {
			ch <- amqpMeta.ErrClosed
			return nil, err
		}
		m.connection = conn
		conn.NotifyClose(ch)
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
	m.registerLock.Lock()
	defer m.registerLock.Unlock()
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
	m.registerLock.Lock()
	defer m.registerLock.Unlock()
	if v, ok := c.(IConsumer); ok {
		m.consumers[c.Identify()] = v
		return c.Link(m)
	}
	if v, ok := c.(IProducer); ok {
		m.producers[c.Identify()] = v
		return c.Link(m)
	}
	return fmt.Errorf("%s is not define", c.Identify())
}

func (m *AMQPManager) NoRelink() {
	m.noAutoRelink <- struct{}{}
}

// AutoRelink auto relink consumer/producer in manager
func (m *AMQPManager) AutoRelink(ctx context.Context) {
	go func() {
		for {
			GetLogger().InfoCtxf(ctx, "begin auto relink")
			time.Sleep(time.Second * 1)
			select {
			case <-m.closeCh:
				for _, c := range m.consumers {
					err := c.Link(m)
					GetLogger().InfoCtxf(ctx, "consumer[%s] retry register err %s", c.Identify(), err)
				}
				for _, c := range m.producers {
					err := c.Link(m)
					GetLogger().InfoCtxf(ctx, "producer[%s] retry register err %s", c.Identify(), err)
				}
			case <-m.noAutoRelink:
				GetLogger().InfoCtxf(ctx, "close auto relink")
				return
			}
		}
	}()
}

// PublishOnce publish message in new connect
func (m *AMQPManager) PublishOnce(ctx context.Context, exchange, router string, encoder Encoder, messages ...IMessage) ([]error, error) {
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

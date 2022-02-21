package amqp

type ICaller interface {
	Identify() string
	Link(m *AMQPManager) error
}

package amqp

type ICaller interface {
	Register(m *AMQPManager) error
}

package amqp

import (
	"context"
	amqpMeta "github.com/streadway/amqp"
	"testing"
	"time"
)

func TestConsume(t *testing.T) {
	ctx := context.Background()

	const (
		TestHost = "192.168.97.224"
		TestPort = int32(5672)
		TestUser = "test"
		TestPWD  = "123456"
		TestPath = "/test"
	)
	m := NewManager(TestHost, TestUser, TestPWD, TestPath, TestPort)
	ShowDefaultLogger(true)
	//_, err := m.GetConnect()
	//if err != nil {
	//	t.Error(err)
	//}

	c := NewConsumer(ctx, "c1", "test.q")
	c.AddDeliveryHandler(func(e amqpMeta.Delivery) {
		t.Log(e.MessageId)
		t.Log(string(e.Body))
		e.Ack(false)
	})
	c.AddCloseHandler(func(e *amqpMeta.Error) {
		t.Error(e.Error())
	})

	err := m.Register(c)
	if err != nil {
		t.Error(err)
	}
	m.AutoRelink(ctx)
	go func() {
		time.Sleep(time.Second * 5)
		m.connection.Close()
		time.Sleep(time.Second * 20)
		m.connection.Close()
		time.Sleep(time.Second * 40)
		m.connection.Close()
	}()

	<-time.After(time.Minute * 5)
}

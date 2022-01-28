package amqp

import (
	"context"
	"testing"
	"time"
)

func TestProducer_Publish(t *testing.T) {
	ctx := context.Background()

	const (
		TestHost = "192.168.97.224"
		TestPort = int32(5672)
		TestUser = "test"
		TestPWD  = "123456"
		TestPath = "/test"
	)
	m := NewManager(TestHost, TestUser, TestPWD, TestPath, TestPort)
	_, err := m.GetConnect()
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(time.Second * 5)
		m.DisConnect()
	}()

	p := &Producer{
		Exchange:      "test.ex",
		closeHandlers: make([]CloseHandler, 0),
		m:             m,
	}

	p.Register(m)

	errs, err := p.Publish(ctx, "test", JSONEncoder{}, &TestMessage{ID: "ticket1"}, &TestMessage{ID: "ticket2"}, &TestMessage{ID: "ticket3"})
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "test", JSONEncoder{}, &TestMessage{ID: "tickasddsdaet1"})
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "test", JSONEncoder{}, &TestMessage{ID: "tickdsaadet1"})
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "test", JSONEncoder{}, &TestMessage{ID: "tickesadt1"})
	println(errs)
	println(err)
}

type TestMessage struct {
	ID string
}

func (m *TestMessage) GetMessageID() string {
	return m.ID
}

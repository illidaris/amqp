package amqp

import (
	"context"
	"testing"
	"time"
)

func TestProducer_Publish(t *testing.T) {
	ctx := context.Background()

	testProducerManager := &AMQPManager{
		TCPSection: TCPSection{
			host: "localhost",
			port: 5672,
			user: "hop",
			pwd:  "123456",
			path: "/hop",
		},
	}

	_, err := testProducerManager.GetConnect()
	if err != nil {
		t.Error(err)
	}
	go func() {
		time.Sleep(time.Second * 5)
		testProducerManager.DisConnect()
	}()

	p := &Producer{
		Exchange:      "diamond.change.ex",
		closeHandlers: make([]CloseHandler, 0),
		m:             testProducerManager,
	}

	p.Register(testProducerManager)

	errs, err := p.Publish(ctx, "diamond.change", JSONEncoder{}, &TestMessage{ID: "ticket1"}, &TestMessage{ID: "ticket2"}, &TestMessage{ID: "ticket3"})
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "diamond.change", JSONEncoder{}, &TestMessage{ID: "tickasddsdaet1"})
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "diamond.change", JSONEncoder{}, &TestMessage{ID: "tickdsaadet1"})
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "diamond.chnge", JSONEncoder{}, &TestMessage{ID: "tickesadt1"})
	println(errs)
	println(err)
}

type TestMessage struct {
	ID string
}

func (m *TestMessage) GetMessageID() string {
	return m.ID
}

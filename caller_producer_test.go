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

	errs, err := p.Publish(ctx, "diamond.change", nil, "ticket1", "ticket2", "ticket3", "ticket4")
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "diamond.change", nil, "xxxx456")
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "diamond.change", nil, "xxxx789")
	println(errs)
	println(err)
	errs, err = p.Publish(ctx, "diamond.chnge", nil, "xx88xx123")
	println(errs)
	println(err)
}

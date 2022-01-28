package amqp

import (
	"testing"
)

func TestAMQPManager_Declare(t *testing.T) {
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

	deadEx := Exchange{
		Element:  Element{"exchange.dlx", true, false, false, nil},
		Kind:     Direct,
		Internal: false,
	}
	deadQ := Queue{
		Element:   Element{"queue.dlx", true, false, false, nil},
		Exclusive: false,
	}
	deadR := Router{
		Ex:        &deadEx,
		Q:         &deadQ,
		Name:      "router.dlx",
		NoWait:    false,
		Arguments: nil,
	}
	args := map[string]interface{}{
		"x-dead-letter-exchange":    "exchange.dlx",
		"x-dead-letter-routing-key": "router.dlx",
	}

	q := Queue{
		Element:   Element{"test.q", true, false, false, args},
		Exclusive: false,
	}
	e := Exchange{
		Element:  Element{"test.ex", true, false, false, nil},
		Kind:     Direct,
		Internal: false,
	}
	r := Router{
		Ex:        &e,
		Q:         &q,
		Name:      "test",
		NoWait:    false,
		Arguments: nil,
	}
	ch, err := m.NewChannel()
	if err != nil {
		t.Error(err)
	}
	defer ch.Close()
	m.Declare(ch, WithQueue(deadQ), WithExchange(deadEx), WithRouter(deadR), WithQueue(q), WithExchange(e), WithRouter(r))
	if err != nil {
		t.Error(err)
	}
}

func ExampleAMQPManager_URL() {
	m := &AMQPManager{
		TCPSection: TCPSection{
			host: "localhost",
			port: 5672,
			user: "test",
			pwd:  "123456",
			path: "/test",
		},
	}
	println(m.URL())
}

//
//func TestAMQPManager_URL(t *testing.T) {
//	const right = "amqp://test:123456@localhost:5672//test"
//	if left := testManager.URL(); left != right {
//		t.Errorf("%s != %s", left, right)
//	}
//}
//
//func ExampleAMQPManager_PublishOnce() {
//	ctx := context.TODO()
//	_, err := testManager.PublishOnce(ctx, "test.ex", "test.router", JSONEncoder{}, &TestMessage2{ID: "test_message_01"})
//	if err != nil {
//		println(err)
//	}
//}
//
//func TestAMQPManager_PublishOnce(t *testing.T) {
//	ctx := context.TODO()
//	_, err := testManager.PublishOnce(ctx, "test.ex", "test.router", JSONEncoder{}, &TestMessage2{ID: "test_message_02"})
//	if err != nil {
//		t.Error(err)
//	}
//}
//
//func TestAmqp(t *testing.T) {
//	_, err := testManager.GetConnect()
//	if err != nil {
//		t.Error(err)
//	}
//	go func() {
//		time.Sleep(time.Second * 60)
//		testManager.DisConnect()
//	}()
//	//NewDeclare()
//	for i := 0; i < 5; i++ {
//		c := NewConsumer("c" + strconv.Itoa(i))
//		err = testManager.Register(c)
//		if err != nil {
//			t.Error(err)
//		}
//	}
//	for {
//		select {
//		case <-time.After(time.Minute * 2):
//			return
//		}
//	}
//}
//
//func NewConsumer(name string) *Consumer {
//	c := &Consumer{
//		Name:      name,
//		Ctx:       context.Background(),
//		QueueName: "diamond.change.q",
//		NoLocal:   false,
//		AutoAck:   false,
//		Exclusive: false,
//		NoWait:    false,
//	}
//	c.DeliveryHandlers = append(c.DeliveryHandlers, func(delivery amqpMeta.Delivery) {
//		//defer func() {
//		//	if delivery.Acknowledger!=nil{
//		//		delivery.Ack(false)
//		//	}
//		//}()
//		var msg string
//		if delivery.Body != nil {
//			msg = string(delivery.Body)
//
//			if strings.Contains(msg, "dead") {
//				delivery.Reject(false)
//			} else {
//				delivery.Ack(false)
//			}
//		}
//
//		println(fmt.Sprintf("%s msg %s", delivery.ConsumerTag, msg))
//	})
//	c.CloseHandlers = append(c.CloseHandlers, func(err *amqpMeta.Error) {
//		println(fmt.Sprintf("%s error %s", c.Name, err))
//	})
//	return c
//}
//
//func NewDeclare(channel *amqpMeta.Channel) {
//	deadEx := Exchange{
//		Element:  Element{"exchange.dlx", true, false, false, nil},
//		Kind:     Direct,
//		Internal: false,
//	}
//	deadQ := Queue{
//		Element:   Element{"queue.dlx", true, false, false, nil},
//		Exclusive: false,
//	}
//	deadR := Router{
//		Ex:        &deadEx,
//		Q:         &deadQ,
//		Name:      "router.dlx",
//		NoWait:    false,
//		Arguments: nil,
//	}
//	args := map[string]interface{}{
//		"x-dead-letter-exchange":    "exchange.dlx",
//		"x-dead-letter-routing-key": "router.dlx",
//	}
//
//	q := Queue{
//		Element:   Element{"diamond.change.q", true, false, false, args},
//		Exclusive: false,
//	}
//	e := Exchange{
//		Element:  Element{"diamond.change.ex", true, false, false, nil},
//		Kind:     Direct,
//		Internal: false,
//	}
//	r := Router{
//		Ex:        &e,
//		Q:         &q,
//		Name:      "diamond.change",
//		NoWait:    false,
//		Arguments: nil,
//	}
//
//	testManager.Declare(channel, WithQueue(deadQ), WithExchange(deadEx), WithRouter(deadR), WithQueue(q), WithExchange(e), WithRouter(r))
//}
//
//type TestMessage2 struct {
//	ID string
//}
//
//func (m *TestMessage2) GetMessageID() string {
//	return m.ID
//}

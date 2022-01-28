package amqp

import (
	"context"
	"errors"
	amqpMeta "github.com/streadway/amqp"
	"time"
)

type DeclareFunc func(channel *amqpMeta.Channel) error
type CloseHandler func(e *amqpMeta.Error)
type DeliveryHandler func(e amqpMeta.Delivery)
type ConfirmHandler func(c amqpMeta.Confirmation)
type ReturnHandler func(r amqpMeta.Return)

func WithExchange(e Exchange) DeclareFunc {
	return e.Declare
}

func WithRouter(e Router) DeclareFunc {
	return e.Bind
}

func WithQueue(e Queue) DeclareFunc {
	return func(channel *amqpMeta.Channel) error {
		mq, err := e.Declare(channel)
		println(mq.Name)
		return err
	}
}

// send message to
func send(ctx context.Context, channel *amqpMeta.Channel, exchange, router string, encoder Encoder, messages ...IMessage) ([]error, error) {
	result := make([]error, 0)
	// on confirm model
	err := channel.Confirm(false)
	if err != nil {
		return result, err
	}
	// listen close channel
	closeCh := channel.NotifyClose(make(chan *amqpMeta.Error, 1))
	// listen confirm channel
	confirmCh := channel.NotifyPublish(make(chan amqpMeta.Confirmation, 1))
	// listen return channel
	returnCh := channel.NotifyReturn(make(chan amqpMeta.Return, 1))

	for _, message := range messages {
		messagePacked, msgErr := packMessage(encoder, message)
		if msgErr != nil {
			result = append(result, msgErr)
			break
		}
		msgErr = channel.Publish(exchange, router, true, false, messagePacked)
		if msgErr != nil {
			result = append(result, msgErr)
			break
		}
		_, msgErr = publishResult(ctx, closeCh, confirmCh, returnCh, time.Second*30)
		if msgErr != nil {
			GetLogger().ErrorCtxf(ctx, "message[%s] publish failed", messagePacked.MessageId)
			result = append(result, msgErr)
			break
		}
		result = append(result, nil)
	}
	return result, nil
}

// publishResult catch result error
func publishResult(ctx context.Context, closeCh <-chan *amqpMeta.Error, configCh <-chan amqpMeta.Confirmation, returnCh <-chan amqpMeta.Return, timeout time.Duration) (*amqpMeta.Return, error) {
	select {
	case c := <-configCh:
		GetLogger().InfoCtxf(ctx, "confirm tag[%d] ack[%t]", c.DeliveryTag, c.Ack)
		if !c.Ack {
			return nil, errors.New("ack is false")
		}
		return nil, nil
	case r := <-returnCh:
		GetLogger().ErrorCtxf(ctx, "return messageId[%s]", r.MessageId)
		return &r, errors.New("delivery return")
	case err := <-closeCh:
		GetLogger().ErrorCtxf(ctx, "channel close %s", err)
		return nil, err
	case <-time.After(timeout):
		GetLogger().ErrorCtxf(ctx, "wait confirm or return timeout[%d]ms", timeout.Milliseconds())
		return nil, errors.New("time out")
	}
}

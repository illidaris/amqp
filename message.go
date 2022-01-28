package amqp

import (
	"github.com/google/uuid"
	amqpMeta "github.com/streadway/amqp"
	"time"
)

type IMessage interface {
	GetMessageID() string
}

func packMessage(encoder Encoder, message IMessage) (amqpMeta.Publishing, error) {
	if encoder == nil {
		encoder = &defaultEncoder{}
	}
	// build message
	messagePacked := amqpMeta.Publishing{
		ContentType:     encoder.GetContentType(),
		ContentEncoding: encoder.GetEncoding(),
		Timestamp:       time.Now(),
		MessageId:       message.GetMessageID(),
	}
	if messagePacked.MessageId == "" {
		messagePacked.MessageId = uuid.NewString()
	}
	data, err := encoder.Encode(message)
	messagePacked.Body = data
	return messagePacked, err
}

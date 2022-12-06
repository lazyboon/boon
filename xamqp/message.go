package xamqp

import (
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
	"time"
)

type Message struct {
	Headers         map[string]interface{}
	ContentType     string
	ContentEncoding string
	DeliveryMode    uint8
	Priority        uint8
	CorrelationId   string
	ReplyTo         string
	Expiration      string
	MessageId       string
	Timestamp       time.Time
	Type            string
	UserId          string
	AppId           string
	Body            []byte
}

func NewMessageDefault(body interface{}) (*Message, error) {
	payload, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(body)
	if err != nil {
		return nil, err
	}
	return &Message{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    uuid.New().String(),
		Timestamp:    time.Now(),
		Body:         payload,
	}, nil
}

func (m *Message) toPublishing() amqp.Publishing {
	return amqp.Publishing{
		Headers:         m.Headers,
		ContentType:     m.ContentType,
		ContentEncoding: m.ContentEncoding,
		DeliveryMode:    m.DeliveryMode,
		Priority:        m.Priority,
		CorrelationId:   m.CorrelationId,
		ReplyTo:         m.ReplyTo,
		Expiration:      m.Expiration,
		MessageId:       m.MessageId,
		Timestamp:       m.Timestamp,
		Type:            m.Type,
		UserId:          m.UserId,
		AppId:           m.AppId,
		Body:            m.Body,
	}
}

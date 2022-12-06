package xamqp

import "github.com/streadway/amqp"

type HandlerType int8

const (
	HandlerTypeAck HandlerType = iota
	HandlerTypeReject
)

type IHandlerResponse interface {
	Type() HandlerType
	Param() bool
}

type handlerAck struct {
	multiple bool
}

func NewHandlerAck(multiple bool) *handlerAck {
	return &handlerAck{multiple: multiple}
}

func (h *handlerAck) Type() HandlerType {
	return HandlerTypeAck
}

func (h *handlerAck) Param() bool {
	return h.multiple
}

type handlerReject struct {
	requeue bool
}

func NewHandlerReject(requeue bool) *handlerReject {
	return &handlerReject{requeue: requeue}
}

func (h *handlerReject) Type() HandlerType {
	return HandlerTypeReject
}

func (h *handlerReject) Param() bool {
	return h.requeue
}

type Handler func(delivery *amqp.Delivery) IHandlerResponse

type Consumer struct {
	Queue     string
	Tag       string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      map[string]interface{}
	Handler   Handler
}

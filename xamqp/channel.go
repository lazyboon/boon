package xamqp

import (
	"github.com/streadway/amqp"
	"sync/atomic"
	"time"
)

type Channel struct {
	*amqp.Channel
	closed int32
}

func NewChannel(connection *amqp.Connection, qos QOS) (*Channel, error) {
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	channel := &Channel{Channel: ch}
	if channel.qos(qos) != nil {
		return nil, err
	}
	go func(connection *amqp.Connection, channel *Channel) {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			if !ok || channel.IsClosed() {
				debug("amqp channel closed")
				_ = channel.Close()
				break
			}
			_ = channel.Close()
			debugf("amqp channel closed, reason: %+v", reason)

			for {
				time.Sleep(1 * time.Second)
				ch, err := connection.Channel()
				if err == nil {
					debug("amqp channel reconnect success")
					channel.Channel = ch
					if channel.qos(qos) != nil {
						continue
					}
					break
				}
			}
			debugf("amqp channel reconnect failed, error: %+v", err)
		}
	}(connection, channel)
	return channel, nil
}

func (c *Channel) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func (c *Channel) Close() error {
	if c.IsClosed() {
		return amqp.ErrClosed
	}
	atomic.StoreInt32(&c.closed, 1)
	return c.Channel.Close()
}

func (c *Channel) Cancel(consumer string, noWait bool) error {
	if c.IsClosed() {
		return amqp.ErrClosed
	}
	atomic.StoreInt32(&c.closed, 1)
	return c.Channel.Cancel(consumer, noWait)
}

func (c *Channel) Consume(consumer *Consumer) <-chan amqp.Delivery {
	deliveries := make(chan amqp.Delivery)
	go func() {
		for {
			d, err := c.Channel.Consume(consumer.Queue, consumer.Tag, consumer.AutoAck, consumer.Exclusive, consumer.NoLocal, consumer.NoWait, consumer.Args)
			if err != nil {
				debugf("amqp consume failed, error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			for msg := range d {
				deliveries <- msg
			}
			time.Sleep(1 * time.Second)
			if c.IsClosed() {
				break
			}
		}
	}()
	return deliveries
}

func (c *Channel) qos(qos QOS) error {
	if qos != (QOS{}) {
		err := c.Channel.Qos(qos.PrefetchCount, qos.PrefetchSize, qos.Global)
		if err != nil {
			debugf("amqp setting channel qos error: %+v", err)
			return err
		}
	}
	return nil
}

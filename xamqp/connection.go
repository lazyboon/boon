package xamqp

import (
	"github.com/streadway/amqp"
	"time"
)

type Connection struct {
	*amqp.Connection
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	connection := &Connection{Connection: conn}
	go func(connection *Connection) {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			if !ok {
				debug("amqp connection closed")
				break
			}
			debugf("amqp connection closed, reason: %v", reason)

			for {
				time.Sleep(1 * time.Second)
				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					debug("amqp reconnect success")
					break
				}
				debugf("amqp reconnect failed, error: %v", err)
			}
		}
	}(connection)
	return connection, nil
}

func (c *Connection) Channel(qos QOS) (*Channel, error) {
	return NewChannel(c.Connection, qos)
}

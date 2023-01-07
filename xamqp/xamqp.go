package xamqp

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
)

var (
	instanceMap   map[string]*AMQP
	DebugCallback func(msg string)
)

func InitWithConfigs(configs []*Config) {
	initInstancesContainer()
	for _, config := range configs {
		instanceMap[config.Alias] = NewAMQP(config)
	}
}

func Connect(alias ...string) *AMQP {
	k := ""
	if len(alias) > 0 {
		k = alias[len(alias)-1]
	}
	return instanceMap[k]
}

func initInstancesContainer() {
	if instanceMap == nil {
		instanceMap = make(map[string]*AMQP)
	}
}

type AMQP struct {
	conf           *Config
	connection     *Connection
	publishChannel *Channel
}

func NewAMQP(conf *Config) *AMQP {
	ans := &AMQP{conf: conf}
	err := ans.dial()
	if err != nil {
		panic(err)
	}
	err = ans.declare()
	if err != nil {
		panic(err)
	}
	return ans
}

type PublishParam struct {
	Exchange  string
	Msg       *Message
	Key       string
	Mandatory bool
	Immediate bool
	Confirm   bool
}

func (a *AMQP) Publish(param *PublishParam) error {
	if param.Exchange == "" {
		return errors.New("amqp publish error: exchange can't empty")
	}
	if param.Msg == nil {
		return errors.New("amqp publish error: msg can't nil")
	}

	// channel
	if a.publishChannel == nil || a.publishChannel.IsClosed() {
		channel, err := a.connection.Channel(QOS{})
		if err != nil {
			return err
		}
		a.publishChannel = channel
	}

	// publish
	err := a.publishChannel.Publish(param.Exchange, param.Key, param.Mandatory, param.Immediate, param.Msg.toPublishing())
	if !param.Confirm {
		return err
	}

	// confirm
	confirmChan := make(chan amqp.Confirmation)
	err = a.publishChannel.Confirm(false)
	if err != nil {
		return err
	}
	a.publishChannel.NotifyPublish(confirmChan)
	for {
		select {
		case d := <-confirmChan:
			if !d.Ack {
				return errors.New("publish rabbitmq message fail")
			}
			return nil
		}
	}
}

func (a *AMQP) Consume(consumer *Consumer, qos QOS) (<-chan amqp.Delivery, error) {
	channel, err := a.connection.Channel(qos)
	if err != nil {
		return nil, err
	}
	c := channel.Consume(consumer)
	return c, nil
}

func (a *AMQP) dial() error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", a.conf.User, a.conf.Password, a.conf.Host, a.conf.Port, a.conf.Vhost)
	connection, err := NewConnection(url)
	if err != nil {
		return err
	}
	a.connection = connection
	return nil
}

func (a *AMQP) declare() error {
	if len(a.conf.Exchanges) == 0 {
		return nil
	}
	// channel
	channel, err := a.connection.Channel(QOS{})
	if err != nil {
		return err
	}
	defer func(channel *Channel) {
		err := channel.Close()
		if err != nil {
			debugf("channel close error: %+v", err)
		}
	}(channel)

	// declare
	err = a.declareExchanges(channel, a.conf.Exchanges)
	if err != nil {
		return err
	}
	return nil
}

func (a *AMQP) declareExchanges(channel *Channel, exchanges []*ExchangeConfig) error {
	for _, exchange := range exchanges {
		// declare exchange
		err := channel.ExchangeDeclare(exchange.Name, exchange.Type, exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, exchange.Args)
		if err != nil {
			return err
		}
		err = a.declareQueues(channel, exchange.Name, exchange.Queues)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AMQP) declareQueues(channel *Channel, exchangeName string, queues []*QueueConfig) error {
	for _, queue := range queues {
		// declare dlx exchange
		_, hasDlxExchange := queue.Args["x-dead-letter-exchange"]
		if queue.DLXExchange != nil && hasDlxExchange {
			err := a.declareExchanges(channel, []*ExchangeConfig{queue.DLXExchange})
			if err != nil {
				return err
			}
		}
		// declare queue
		_, err := channel.QueueDeclare(queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, queue.NoWait, queue.Args)
		if err != nil {
			return err
		}
		// bind
		fmt.Println(queue.Name, queue.RoutingKey, exchangeName)
		err = channel.QueueBind(queue.Name, queue.RoutingKey, exchangeName, false, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func debug(msg string) {
	if DebugCallback != nil {
		DebugCallback(msg)
	}
}

func debugf(format string, args ...interface{}) {
	if DebugCallback != nil {
		DebugCallback(fmt.Sprintf(format, args...))
	}
}

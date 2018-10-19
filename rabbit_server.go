package core

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"strconv"
	"time"
)

var rabbit_server_stor = map[string]*RabbitServer{}

type RabbitPort struct {
	name     string
	channel  *amqp.Channel
	queue    string
	exchange string
}

type Publisher RabbitPort

func NewPublisher(name string, rs *RabbitServer) *Publisher {
	INFO("init exch: " + name)
	p := &Publisher{
		name:     name,
		channel:  rs.GetChannel(),
		queue:    "",
		exchange: name,
	}

	err := p.channel.ExchangeDeclare(
		name,
		"fanout",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		panic(err)
	}

	return p
}

func (p *Publisher) Send(msg interface{}) (err error) {
	var bmsg []byte

	switch tmsg := msg.(type) {
	case string:
		bmsg = []byte(tmsg)
	case []byte:
		bmsg = tmsg
	}

	err = p.channel.Publish(
		p.exchange,
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bmsg,
		},
	)
	return err
}

func (p *Publisher) SendSequence(args ...interface{}) (err error) {
	pack := MsgPack(args)
	return p.Send(pack)
}

type Consumer RabbitPort

func NewConsumer(name string, rs *RabbitServer, exclusive bool) *Consumer {
	INFO("init consumer: " + name)
	c := &Consumer{
		name:     name,
		channel:  rs.GetChannel(),
		exchange: "",
	}

	if exclusive {
		name = ""
	}

	c.channel.NotifyClose(rs.cerr)

	queue, err := c.channel.QueueDeclare(
		name,
		true,
		false,
		exclusive,
		false,
		nil)

	if err != nil {
		panic(err)
	}

	c.queue = queue.Name

	return c
}

func (c *Consumer) Listen(callback func([]interface{})) {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		ERROR(fmt.Sprintf("Consumer: Listen: Consume: queue: %v, reason: %v", c.name, err.Error()))
	}

	go func() {
		for message := range msgs {
			args := []interface{}{}

			err := MsgUnpack(message.Body, &args)
			if err != nil {
				ERROR(fmt.Sprintf("Consumer: Listen: Unpack: queue: %v, message: %v, reason: %v", c.name, ToString(message.Body), err.Error()))
				continue
			}

			callback(args)
		}
	}()
}

func (c *Consumer) RPCListen(callback_map map[string]func([]interface{}) interface{}) {
	c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		ERROR(fmt.Sprintf("Consumer: RPCListen: Consume: queue: %v, reason: %v", c.name, err.Error()))
	}

	go func() {
		for message := range msgs {
			args := []interface{}{}

			err := MsgUnpack(message.Body, &args)
			if err != nil {
				ERROR(fmt.Sprintf("Consumer: RPCListen: Unpack: queue: %v, message: %v, reason: %v", c.name, ToString(message.Body), err.Error()))
				continue
			}

			meth_number := ToString(args[0])
			var res interface{}

			if method, ok := callback_map[meth_number]; ok {
				if len(args) == 1 {
					res = method([]interface{}{})
				} else {
					res = method(args[1:])
				}

				byte_res := MsgPack(res)

				err = c.channel.Publish(
					"",              // exch
					message.ReplyTo, // routing key
					false,           // mandatory
					false,           // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: message.CorrelationId,
						Body:          byte_res,
					})
				if err != nil {
					ERROR(fmt.Sprintf("Consumer: RPCListen: publish result: Failed to publish a message: err: %v", err.Error()))
					continue
				}

				message.Ack(false)

			} else {
				message := fmt.Sprintf("rabbit_server: RPCListen: method with number %v not found", meth_number)
				ERROR(message)
				continue
			}
		}
	}()
}

type RabbitServer struct {
	name              string
	rabbit_connection *amqp.Connection
	publishers        map[string]*Publisher
	consumers         map[string]*Consumer
	cerr              chan *amqp.Error
}

func NewRabbitServer(srv_name string, conn string) (*RabbitServer, error) {
	connection, err := amqp.Dial(conn)
	if err != nil {
		return nil, err
	}

	c := make(chan *amqp.Error, 2)
	connection.NotifyClose(c)

	go func(tt chan *amqp.Error) {
		<-tt
		time.Sleep(time.Second)
		os.Exit(13)
	}(c)

	serv := &RabbitServer{
		name:              srv_name,
		rabbit_connection: connection,
		publishers:        map[string]*Publisher{},
		consumers:         map[string]*Consumer{},
		cerr:              c,
	}

	return serv, nil
}

func SetRabbitServer(name, conn string) (*RabbitServer, error) {
	rs, err := NewRabbitServer(name, conn)
	if err != nil {
		return nil, err
	}
	rabbit_server_stor[name] = rs
	return rs, nil
}

func GetRabbitServer(name string) (*RabbitServer, error) {
	if rs, ok := rabbit_server_stor[name]; ok {
		return rs, nil
	}

	time.Sleep(time.Second * 5)

	panic(fmt.Sprintf("no rabbit server with name [%v]", name))
}

func (rs *RabbitServer) GetChannel() *amqp.Channel {
	channel, err := rs.rabbit_connection.Channel()
	if err != nil {
		panic(err)
	}
	return channel
}

func (rs *RabbitServer) GetPublisher(name string) *Publisher {
	if p, ok := rs.publishers[name]; ok {
		return p
	} else {
		p := NewPublisher(name, rs)
		rs.publishers[name] = p
		return p
	}
}

func (rs *RabbitServer) GetConsumer(name string) *Consumer {
	if c, ok := rs.consumers[name]; ok {
		return c
	} else {
		c := NewConsumer(name, rs, false)
		rs.consumers[name] = c
		return c
	}
}

func (rs *RabbitServer) GetExclusiveConsumer(name string) *Consumer {
	if c, ok := rs.consumers[name]; ok {
		return c
	} else {
		c := NewConsumer(name, rs, true)
		rs.consumers[name] = c
		return c
	}
}
func (rs *RabbitServer) RPCSend(service string, method string, args ...interface{}) (res interface{}, err error) {
	pack := []interface{}{method}
	pack = append(pack, args...)
	packb := MsgPack(pack)

	cons_name := method + "_rpc"
	consumer := rs.GetExclusiveConsumer(cons_name)

	msgs, err := consumer.channel.Consume(
		consumer.queue, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		msg := fmt.Sprintf("RabbitServer: RPCSend: consumer.Consume: queue: %v, err: %v", cons_name, err.Error())
		ERROR(msg)
		return res, errors.New(msg)
	}

	t := time.Now()
	tu := int(t.Unix())
	corrId := strconv.Itoa(tu)

	err = consumer.channel.Publish(
		"",              // exch
		service+"_main", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       consumer.queue,
			Body:          packb,
		})
	if err != nil {
		msg := fmt.Sprintf("RabbitServer: RPCSend: consumer.Publish: queue: %v, err: %v", cons_name, err.Error())
		ERROR(msg)
		return res, errors.New(msg)
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			MsgUnpackScalar(d.Body, &res)
			break
		}
	}

	return res, nil
}

func (rs *RabbitServer) BindConsumer(queue string, exchange string, temporary bool) (*Consumer, error) {
	var pQueue *Consumer
	if temporary {
		pQueue = rs.GetExclusiveConsumer(exchange + "_tmp")
	} else {
		pQueue = rs.GetConsumer(queue)
	}

	pExchange := rs.GetPublisher(exchange)

	err := pExchange.channel.QueueBind(
		pQueue.queue,       // queue name
		"",                 // routing key
		pExchange.exchange, // exch
		false,
		nil,
	)
	if err != nil {
		msg := fmt.Sprintf("RabbitServer: BindFanout: QueueBind: queue: %v, exch: %v, reason: %v", queue, exchange, err.Error())
		ERROR(msg)
		return nil, errors.New(msg)
	}

	return pQueue, nil
}

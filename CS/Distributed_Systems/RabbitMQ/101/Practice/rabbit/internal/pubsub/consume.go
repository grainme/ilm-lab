package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RabbitUrl = "amqp://guest:guest@localhost:5672/"
)

type AckType string

const (
	Ack         AckType = "ack"
	NackRequeue AckType = "nackRequeue"
	NackDiscard AckType = "nackDiscard"
)

type SimpleQueueType string

const (
	Durable   SimpleQueueType = "durable"
	Transient SimpleQueueType = "transient"
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create channel: %v", err)
	}

	isExclusive, isAutoDelete, isDurable := false, false, false
	if queueType == Transient {
		isAutoDelete, isExclusive = true, true
	}
	if queueType == Durable {
		isDurable = true
	}

	amqpQueue, err := ch.QueueDeclare(queueName, isDurable, isAutoDelete, isExclusive, false, amqp.Table{"x-dead-letter-exchange": "peril_dlx"})
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %v", err)
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue: %v", err)
	}

	return ch, amqpQueue, nil
}

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
	unmarshaller := func(b []byte) (T, error) {
		var deliveryBody T
		err := json.Unmarshal(b, &deliveryBody)

		return deliveryBody, err
	}

	return subscribe(conn, exchange, queueName, key, queueType, handler, unmarshaller)
}

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {

	unmarshaller := func(b []byte) (T, error) {
		var deliveryBody T
		buffer := bytes.NewBuffer(b)

		decoder := gob.NewDecoder(buffer)
		err := decoder.Decode(&deliveryBody)

		return deliveryBody, err
	}

	return subscribe(conn, exchange, queueName, key, queueType, handler, unmarshaller)
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		return fmt.Errorf("could not set QoS: %v", err)
	}
	deliveryChan, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveryChan {
			deliveryBody, err := unmarshaller(delivery.Body)
			if err != nil {
				log.Fatalf("could not decode delivery body: %v", err)
			}

			ackType := handler(deliveryBody)
			switch ackType {
			case Ack:
				log.Println("Delivery acknowledged")
				delivery.Ack(false)
			case NackRequeue:
				log.Println("Delivery negeative acknowledged and requeued")
				delivery.Nack(false, true)
			case NackDiscard:
				log.Println("Delivery negative acknowledged and discraded")
				delivery.Nack(false, false)
			}
		}
	}()

	return nil
}

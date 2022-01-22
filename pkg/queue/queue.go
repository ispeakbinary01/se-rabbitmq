package queue

import (
	"fmt"
	"log"

	"github.com/ispeakbinary01/se-rabbitmq/pkg/conn"
	"github.com/streadway/amqp"
)

func NewMessageQueues(names []string) (queues []*amqp.Queue, err error) {
	conn, err := conn.RabbitMQConnection()
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}
	defer ch.Close()
	for _, name := range names {
		q, err := ch.QueueDeclare(name, false, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		err = ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("Hello %s", name)),
		})
		if err != nil {
			return nil, err
		}
		queues = append(queues, &q)
	}
	fmt.Println("Queues created: ", queues)
	return queues, nil
}

func ProduceToQueue(name string, body []byte) error {
	conn, err := conn.RabbitMQConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(name, false, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}
	return nil
}

func ConsumeFromQueue(name string) error {
	conn, err := conn.RabbitMQConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("consuming from queue %s and the message is %d", name, d.Body)
		}
	}()
	<-forever
	return nil
}

package queue

import (
	"fmt"

	"github.com/ispeakbinary01/se-rabbitmq/conn"
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

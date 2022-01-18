package conn

import (
	"fmt"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func RabbitMQConnection() (*amqp.Connection, error) {
	if conn == nil {
		var err error
		fmt.Println("Connecting to local server")
		conn, err = amqp.Dial("amqp://guest:guest@localhost:5672")
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("Connected and returning connection")
	return conn, nil
}

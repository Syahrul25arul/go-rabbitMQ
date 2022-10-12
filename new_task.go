package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failAnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func main() {
	fmt.Println("RABBIT MQ TUTORIAL")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failAnError(err, "Failed To Connect RabbitMQ")
	defer conn.Close()

	fmt.Println("successfull connect rabbit mq")

	ch, err := conn.Channel()
	failAnError(err, "Failed To Open Channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Test Queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failAnError(err, "Failed to declare queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)

	err = ch.PublishWithContext(
		ctx,    // context with timeout
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	failAnError(err, "Failed to publish message")

	log.Printf(" [x] Sent %s\n", body)
}

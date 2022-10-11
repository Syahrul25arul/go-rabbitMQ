package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failAnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
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

	body := "Hello World"

	err = ch.PublishWithContext(
		ctx,    // context with timeout
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failAnError(err, "Failed to publish message")

	log.Printf(" [x] Sent %s\n", body)
}

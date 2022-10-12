package main

import (
	"bytes"
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
	failAnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Test Queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failAnError(err, "Failed to declare a queue")

	msg, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no local
		false,  // no-wait
		nil,    // arguments
	)
	failAnError(err, "Failed to register a consumer")

	forever := make(chan struct{})

	go func() {
		for d := range msg {
			log.Printf("Recived a message %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

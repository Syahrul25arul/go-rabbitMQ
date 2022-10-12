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
		"task_queue_new", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failAnError(err, "Failed to declare a queue")

	// ini setingan untuk mengatur consumer bahwa service ini tidak akan menerima data dari server rabbitMQ jika data nya masih ada dan diproses.
	// maka server rabbitMQ akan secara otomatis mengirim data service lain dengan queue yang sama dengan ini
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failAnError(err, "Failed to set QoS")

	// jika auto-ack di set false, maka jika consumer penerima data mati dan data dari rabbitmq belum di selesai proses
	// semua data baik yang dikirim atau yang berada di queue akan dipulihkan kembali
	msg, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		false,  // auto-ack
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
			// time.Sleep(20 * time.Second)
			time.Sleep(t * time.Second)
			log.Printf("Done")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

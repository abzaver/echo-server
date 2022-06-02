package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Connecting to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declaring the queue (ensure it matches the producer's configuration)
	queue, err := ch.QueueDeclare(
		"example_queue", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Registering the consumer
	msgs, err := ch.Consume(
		queue.Name, // queue name
		"",         // consumer tag
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	// Creating a channel to wait for incoming messages
	forever := make(chan bool)

	// Processing the received messages
	go func() {
		for d := range msgs {
			switch d.ContentType {
			case "application/json":
				handleJSONMessage(d.Body)
			case "application/xml":
				handleXMLMessage(d.Body)
			case "text/plain":
				handleTextMessage(d.Body)
			case "application/octet-stream":
				handleByteArrayMessage(d.Body)
			default:
				log.Printf("Received unknown content type: %s", d.ContentType)
			}
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}

// handleJSONMessage processes JSON messages
func handleJSONMessage(body []byte) {
	var message map[string]string
	err := json.Unmarshal(body, &message)
	if err != nil {
		log.Printf("Failed to unmarshal JSON: %s", err)
		return
	}
	log.Printf("Received JSON message: %v", message)
}

// handleXMLMessage processes XML messages
func handleXMLMessage(body []byte) {
	xmlData := string(body)
	log.Printf("Received XML message: %s", xmlData)
}

// handleTextMessage processes plain text messages
func handleTextMessage(body []byte) {
	text := string(body)
	log.Printf("Received text message: %s", text)
}

// handleByteArrayMessage processes byte array messages
func handleByteArrayMessage(body []byte) {
	log.Printf("Received byte array message: %v", body)
}

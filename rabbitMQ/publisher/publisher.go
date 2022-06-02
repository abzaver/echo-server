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

	// Declaring a queue
	queue, err := ch.QueueDeclare(
		"example_queue", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Sending a JSON message
	jsonData, err := json.Marshal(map[string]string{"key": "value"})
	failOnError(err, "Failed to marshal JSON")
	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		})
	failOnError(err, "Failed to publish JSON message")

	// Sending an XML message
	xmlData := `<example><key>value</key></example>`
	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/xml",
			Body:        []byte(xmlData),
		})
	failOnError(err, "Failed to publish XML message")

	// Sending a plain text message
	textData := "This is a plain text message"
	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(textData),
		})
	failOnError(err, "Failed to publish text message")

	// Sending a byte array message
	byteArrayData := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        byteArrayData,
		})
	failOnError(err, "Failed to publish byte array message")

	log.Println("All messages were successfully published")
}

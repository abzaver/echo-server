package main

import (
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

	// Declaring Direct, Fanout, and Topic Exchanges
	err = ch.ExchangeDeclare("direct_exchange", "direct", true, false, false, false, nil)
	failOnError(err, "Failed to declare a direct exchange")

	err = ch.ExchangeDeclare("fanout_exchange", "fanout", true, false, false, false, nil)
	failOnError(err, "Failed to declare a fanout exchange")

	err = ch.ExchangeDeclare("topic_exchange", "topic", true, false, false, false, nil)
	failOnError(err, "Failed to declare a topic exchange")

	// Declaring Queues and Binding them to Exchanges
	_, err = ch.QueueDeclare("direct_queue", true, false, false, false, nil)
	failOnError(err, "Failed to declare a direct queue")

	_, err = ch.QueueDeclare("fanout_queue", true, false, false, false, nil)
	failOnError(err, "Failed to declare a fanout queue")

	_, err = ch.QueueDeclare("topic_queue", true, false, false, false, nil)
	failOnError(err, "Failed to declare a topic queue")

	err = ch.QueueBind("direct_queue", "direct_key", "direct_exchange", false, nil)
	failOnError(err, "Failed to bind direct queue to exchange")

	err = ch.QueueBind("fanout_queue", "", "fanout_exchange", false, nil)
	failOnError(err, "Failed to bind fanout queue to exchange")

	err = ch.QueueBind("topic_queue", "topic.*", "topic_exchange", false, nil)
	failOnError(err, "Failed to bind topic queue to exchange")

	// Consuming messages from each queue
	consumeMessages(ch, "direct_queue")
	consumeMessages(ch, "fanout_queue")
	consumeMessages(ch, "topic_queue")

	log.Println("Server is running and waiting for messages. To exit press CTRL+C")
	select {} // Infinite loop to keep the server running
}

// consumeMessages listens for incoming messages from a specified queue
func consumeMessages(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message from %s: %s", queueName, d.Body)
		}
	}()
}

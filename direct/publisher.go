package main

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishMessage() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ.")
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel.")
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", //exchange name
		"direct",      //exchange type
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"logs_direct", //exchange name
		"info",     //routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		failOnError(err, "Failed to publish a message.")
		return
	}
	fmt.Println("Message Published:", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}

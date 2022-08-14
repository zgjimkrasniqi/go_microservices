package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// Try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer func(rabbitConn *amqp.Connection) {
		err := rabbitConn.Close()
		if err != nil {
			log.Println(err)
		}
	}(rabbitConn)

	// Start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// Create a consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// Watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	/*
		RabbitMQ sometimes can take a while to startup, so a go routine is needed
		that will attempt to connect a fixed number of times
	*/

	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// Do not continue until rabbit is ready
	// amqp.Dial(amqp://username:password@host")
	for {
		// If you want to test it locally, but since we are using docker
		// localhost should be replaced with rabbitmq (as we have specified in the docker compose)
		// c, err := amqp.Dial("amqp://guest:guest@localhost")
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			// If there is an error, we cannot connect right now
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ")
			break
		}

		if counts > 5 {
			// If we cannot connect after 5 times, something is wrong
			fmt.Println("Something is wrong: ", err)
			return nil, err
		}

		// For every time that we tried and did not receive a connection
		// we will backoff, and this backoff time will increase every time
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

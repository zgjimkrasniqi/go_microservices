package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

// Consumer Type that is going to be used to receive events from the queue
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	// We are going to have to open up a channel and declare an exchange
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// Payload Type that is going to be used to push events to the queue
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	// Go to receiverConsumer.connection.channel, get something from it
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			return
		}
	}(ch)

	// Since there is no error, the channel is initialized
	// We need to get a Logger Queue, declare a random queue and use it
	q, err := declareLoggerQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"authentication_logger_exchange",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exhange, Queue] [authentication_logger_exchange, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// authenticate

	default:
		// Do something for default case
	}
}

func logEvent(entry Payload) error {
	// Create some json and send to the logger microservice
	jsonData, _ := json.Marshal(entry)

	// Call the service
	// url: http://<name of the service specified in docker-compose.yml>/...
	request, err := http.NewRequest("POST", "http://logger-service/insert_log", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

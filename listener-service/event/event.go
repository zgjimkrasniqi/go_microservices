package event

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"authentication_logger_exchange",
		"topic",
		true,
		false,
		false, // false since we are going to use it between microservices
		false,
		nil,
	)
}

func declareLoggerQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"loggerQueue",
		false,
		false,
		true,
		false,
		nil,
	)
}

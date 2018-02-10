package main

import "github.com/streadway/amqp"

func main() {
	// Connect to rabbitmq
	uri := amqp.URI{
		Scheme:   "amqp",
		Host:     MQSERVERNAME,
		Port:     5672,
		Username: "user",
		Password: "kQS5MZHEFC",
		Vhost:    "/",
	}
	conn, err := amqp.Dial(uri.String())
	failOnError(err, "Failed to connect to RabbitMQ")

	newsiteChannel := make(chan NewSiteEvent)
	forever := make(chan bool)
	dequeueEvents(conn, newsiteChannel)
	consumeEvents(newsiteChannel)
	<-forever
}

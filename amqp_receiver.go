package main

import (
	"github.com/streadway/amqp"
	"encoding/json"
	"log"
)

func dequeueEvents(conn *amqp.Connection, newsiteChannel chan NewSiteEvent) {
	ch, err := conn.Channel()
	failOnError(err,"Failed to open AMQP channel")
	newsiteQ, _ := ch.QueueDeclare(
		newsiteQueueName, 		// name
		false,          // durable
		false,		// delete when usused
		false,			// exclusive
		false,			// no-wait
		nil,				// arguments
	)
	newsiteIn, _ := ch.Consume(
		newsiteQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	go func() {
		for {
			select {
			case newsiteRaw := <-newsiteIn:
				dispatchNewSite(newsiteRaw, newsiteChannel)
			}
		}
	}()
}

func dispatchNewSite(newsiteRaw amqp.Delivery, out chan NewSiteEvent) {
	var event NewSiteEvent
	err := json.Unmarshal(newsiteRaw.Body, &event)
	if err == nil {
		out <- event
	} else {
		log.Printf("Failed to deserialize raw newsite event %s from queue: %v", newsiteRaw.Body, err)
	}
}

package main

import (
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const (
	MQSERVERNAME = "historical-mandrill-rabbitmq"
	newsiteQueueName = "new_site"
	dataDir = "/data"
	sitesDir = "/data/sites"
)

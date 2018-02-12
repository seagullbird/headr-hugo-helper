package main

import (
	"github.com/seagullbird/headr-common/mq_helper"
)

func main() {
	receiver := mq_helper.NewReceiver()
	receiver.RegisterListener(newsiteQueueName, generateNewSite)

	forever := make(chan bool)
	<-forever
}

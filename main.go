package main

import (
	"github.com/go-kit/kit/log"
	"github.com/seagullbird/headr-common/mq/client"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/seagullbird/headr-hugo-helper/config"
	"github.com/seagullbird/headr-hugo-helper/listeners"
	"os"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		servername = os.Getenv("RABBITMQ_SERVER")
		username   = os.Getenv("RABBITMQ_USER")
		passwd     = os.Getenv("RABBITMQ_PASS")
	)

	receiver, err := receive.NewReceiver(client.New(servername, username, passwd), logger)
	if err != nil {
		logger.Log("error_desc", "receive.NewReceiver failed", "error", err)
		return
	}
	receiver.RegisterListener(config.NewsiteQueueName, listeners.MakeGenerateNewSiteListener(logger))
	receiver.RegisterListener(config.RegenerateQueueName, listeners.MakeReGenerateListener(logger))

	forever := make(chan bool)
	<-forever
}

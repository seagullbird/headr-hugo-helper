package main

import (
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/client"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/go-kit/kit/log"
	"os"
	"github.com/seagullbird/headr-hugo-helper/listeners"
	"github.com/seagullbird/headr-hugo-helper/config"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		servername = mq.MQSERVERNAME
		username   = mq.MQUSERNAME
		passwd     = mq.MQSERVERPWD
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

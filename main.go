package main

import (
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/go-kit/kit/log"
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
		servername = mq.MQSERVERNAME
		username   = mq.MQUSERNAME
		passwd     = mq.MQSERVERPWD
	)

	conn, err := mq.MakeConn(servername, username, passwd)
	if err != nil {
		logger.Log("error_desc", "mq.MakeConn failed", "error", err)
		return
	}
	receiver, err := receive.NewReceiver(conn, logger)
	if err != nil {
		logger.Log("error_desc", "receive.NewReceiver failed", "error", err)
		return
	}
	receiver.RegisterListener(newsiteQueueName, makegenerateNewSiteListener(logger))

	forever := make(chan bool)
	<-forever
}

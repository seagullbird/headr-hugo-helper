package main

import (
	"github.com/go-kit/kit/log"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/streadway/amqp"
	"encoding/json"
	"fmt"
	"github.com/seagullbird/headr-common/mq"
)


func makegenerateNewSiteListener(logger log.Logger) receive.Listener {
	return func(delivery amqp.Delivery) {
		var event mq.NewSiteEvent
		err := json.Unmarshal(delivery.Body, &event)
		if err != nil {
			logger.Log("error_desc", "Failed to unmarshal event","error", err, "raw-message:", delivery.Body)
			return
		}

		logger.Log("info", "Received newsite event", "event", event)
		sitepath := filepath.Join(sitesDir, event.Email, event.SiteName)

		if _, err := os.Stat(sitepath); err == nil || !os.IsNotExist(err) {
			logger.Log("info", fmt.Sprintf("Path %s already exists.", sitepath))
			return
		}
		cmd := exec.Command("hugo", "new", "site", sitepath)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			logger.Log("error_desc",  "Failed to generate new site", "error", err)
		}
	}
}
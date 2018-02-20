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


func makeGenerateNewSiteListener(logger log.Logger) receive.Listener {
	return func(delivery amqp.Delivery) {
		var event mq.NewSiteEvent
		err := json.Unmarshal(delivery.Body, &event)
		if err != nil {
			logger.Log("error_desc", "Failed to unmarshal event","error", err, "raw-message:", delivery.Body)
			return
		}

		logger.Log("info", "Received newsite event", "event", event)
		sitepath := filepath.Join(sitesDir, event.Email, event.SiteName)
		siteSourcePath := filepath.Join(sitepath, "source")
		sitePublicPath := filepath.Join(sitepath, "public")

		if _, err := os.Stat(sitepath); err == nil || !os.IsNotExist(err) {
			logger.Log("info", fmt.Sprintf("Path %s already exists.", sitepath))
			return
		}
		if err := runCommand("hugo", "new", "site", siteSourcePath); err != nil {
			logger.Log("error_desc",  "failed to generate new site source", "error", err)
			return
		}
		if err := runCommand(
			"hugo",
			"--source", siteSourcePath,
			"--destination", sitePublicPath,
			"--themesDir", themesDir,
			"--theme", initialThemeName,
			"--config", filepath.Join(configsDir, initialThemeName, "config.toml")); err != nil {
			logger.Log("error_desc",  "failed to generate new site public", "error", err)
			return
		}
	}
}

func runCommand(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
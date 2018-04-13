package listeners

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/seagullbird/headr-hugo-helper/config"
	"github.com/streadway/amqp"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// MakeGenerateNewSiteListener returns a Listener listening to new_site queue
func MakeGenerateNewSiteListener(logger log.Logger) receive.Listener {
	return func(delivery amqp.Delivery) {
		var event mq.SiteUpdatedEvent
		err := json.Unmarshal(delivery.Body, &event)
		if err != nil {
			logger.Log("error_desc", "Failed to unmarshal event", "error", err, "raw-message:", delivery.Body)
			return
		}

		logger.Log("info", "Received newsite event", "event", event)
		sitepath := filepath.Join(config.SitesDir, strconv.Itoa(int(event.SiteID)))
		siteSourcePath := filepath.Join(sitepath, "source")
		sitePublicPath := filepath.Join(sitepath, "public")
		publicConfigPath := filepath.Join(config.ConfigsDir, event.Theme, "config.json")

		if _, err := os.Stat(sitepath); err == nil || !os.IsNotExist(err) {
			logger.Log("info", fmt.Sprintf("Path %s already exists.", sitepath))
			return
		}
		if err := runCommand("hugo", "new", "site", siteSourcePath, "--format", "json"); err != nil {
			logger.Log("error_desc", "failed to generate new site source", "error", err)
			return
		}
		if err := runCommand("cp", publicConfigPath, siteSourcePath); err != nil {
			logger.Log("error_desc", "failed to cp site config", "error", err)
			return
		}
		baseURL := "/" + strconv.Itoa(int(event.SiteID)) + "/"
		if err := reGenerate(baseURL, siteSourcePath, sitePublicPath); err != nil {
			logger.Log("error_desc", "failed to generate new site public", "error", err)
			return
		}
	}
}

// MakeReGenerateListener returns a Listener listening to re_generate queue
func MakeReGenerateListener(logger log.Logger) receive.Listener {
	return func(delivery amqp.Delivery) {
		var event mq.SiteUpdatedEvent
		err := json.Unmarshal(delivery.Body, &event)
		if err != nil {
			logger.Log("error_desc", "Failed to unmarshal event", "error", err, "raw-message:", delivery.Body)
			return
		}

		logger.Log("info", "Received regenerate event", "event", event)
		sitepath := filepath.Join(config.SitesDir, strconv.Itoa(int(event.SiteID)))
		siteSourcePath := filepath.Join(sitepath, "source")
		sitePublicPath := filepath.Join(sitepath, "public")
		baseURL := "/" + strconv.Itoa(int(event.SiteID)) + "/"
		if err := reGenerate(baseURL, siteSourcePath, sitePublicPath); err != nil {
			logger.Log("error_desc", "failed to re-generate site", "error", err)
			return
		}
	}
}

func reGenerate(baseURL, source, destination string) error {
	if config.Dev == "true" {
		return runCommand("hugo", "--source", source, "--destination", destination)
	}
	return runCommand("hugo", "--baseURL", baseURL, "--source", source, "--destination", destination)
}

func runCommand(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

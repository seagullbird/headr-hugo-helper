package listeners

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
	"github.com/seagullbird/headr-hugo-helper/config"
	"strconv"
)

// MakeGenerateNewSiteListener returns a Listener listening to new_site queue
func MakeGenerateNewSiteListener(logger log.Logger) receive.Listener {
	return func(delivery amqp.Delivery) {
		var event mq.SiteUpdatedEvent
		err := json.Unmarshal(delivery.Body, &event)
		if err != nil {
			logger.Log("error_desc", "Failed to unmarshal event","error", err, "raw-message:", delivery.Body)
			return
		}

		logger.Log("info", "Received newsite event", "event", event)
		sitepath := filepath.Join(config.SitesDir, strconv.Itoa(int(event.SiteId)))
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
		if err := reGenerate(siteSourcePath, sitePublicPath, event.Theme); err != nil {
			logger.Log("error_desc",  "failed to generate new site public", "error", err)
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
			logger.Log("error_desc", "Failed to unmarshal event","error", err, "raw-message:", delivery.Body)
			return
		}

		logger.Log("info", "Received regenerate event", "event", event)
		sitepath := filepath.Join(config.SitesDir, strconv.Itoa(int(event.SiteId)))
		siteSourcePath := filepath.Join(sitepath, "source")
		sitePublicPath := filepath.Join(sitepath, "public")
		if err := reGenerate(siteSourcePath, sitePublicPath, event.Theme); err != nil {
			logger.Log("error_desc",  "failed to re-generate site", "error", err)
			return
		}
	}
}

func reGenerate(source, destination, theme string) error {
	return runCommand(
		"hugo",
		"--source", source,
		"--destination", destination,
		"--themesDir", config.ThemesDir,
		"--theme", theme,
		"--config", filepath.Join(config.ConfigsDir, theme, "config.toml"))
}

func runCommand(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
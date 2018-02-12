package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/seagullbird/headr-common/mq_helper"
	"github.com/streadway/amqp"
	"encoding/json"
	"fmt"
)

func generateNewSite(delivery amqp.Delivery) {
	var event mq_helper.NewSiteEvent
	err := json.Unmarshal(delivery.Body, &event)
	if err != nil {
		log.Println("Failed to unmarshal event:", err, "raw message:", delivery.Body)
		return
	}

	log.Println("Received newsite event:", event)
	sitepath := filepath.Join(sitesDir, event.Email, event.SiteName)

	if _, err := os.Stat(sitepath); err == nil || !os.IsNotExist(err) {
		log.Println(fmt.Sprintf("Path %s already exists.", sitepath))
		return
	}
	cmd := exec.Command("hugo", "new", "site", sitepath)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println( "Failed to generate new site", err)
	}
}
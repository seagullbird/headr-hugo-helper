package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func consumeEvents(newsiteChannel chan NewSiteEvent) {
	go func() {
		log.Println("Started event consumer goroutine")
		for {
			select {
			case newsite := <-newsiteChannel:
				generateNewSite(newsite)
			}
		}

	}()
}

func generateNewSite(event NewSiteEvent) {
	log.Printf("Received newsite event: %s", event)
	sitepath := filepath.Join(sitesDir, event.Email, event.SiteName)

	if _, err := os.Stat(sitepath); err == nil || !os.IsNotExist(err) {
		log.Printf("Path %s already exists: %v", sitepath, err)
		return
	}
	cmd := exec.Command("hugo", "new", "site", sitepath)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		failOnError(err, "Failed to generate new site")
	} else {
		log.Println("Successfully generated new site.")
	}
}

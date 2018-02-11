package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"os"
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
	cmd := exec.Command("hugo", "new", "site", filepath.Join(sitesDir, event.Email, event.SiteName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		failOnError(err, "Failed to generate new site")
	}
}

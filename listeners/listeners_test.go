package listeners_test

import (
	"testing"
	"os"
	"github.com/go-kit/kit/log"
	"github.com/seagullbird/headr-common/mq"
	"github.com/seagullbird/headr-common/mq/client"
	"github.com/seagullbird/headr-common/mq/dispatch"
	"time"
	"github.com/seagullbird/headr-common/mq/receive"
	"github.com/seagullbird/headr-hugo-helper/listeners"
	"github.com/seagullbird/headr-hugo-helper/config"
	"path/filepath"
	"strconv"
)

func TestListeners(t *testing.T) {
	// Logging
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// rabbitmq server
	var (
		servername = os.Getenv("RABBITMQ_PORT_5672_TCP_ADDR")
		username   = "guest"
		passwd     = "guest"
	)

	// New dispatcher
	dispatcher, err := dispatch.NewDispatcher(client.New(servername, username, passwd), logger)
	if err != nil {
		panic(err)
	}

	// New receiver
	receiver, err := receive.NewReceiver(client.New(servername, username, passwd), logger)
	if err != nil {
		t.Fatal(err)
	}
	receiver.RegisterListener(config.NewsiteQueueName, listeners.MakeGenerateNewSiteListener(logger))
	receiver.RegisterListener(config.RegenerateQueueName, listeners.MakeReGenerateListener(logger))


	// Dispatch a NewSite Message
	fakeSiteId := 123
	fakeTheme := "test_theme"
	msg := mq.SiteUpdatedEvent{
		SiteID: uint(fakeSiteId),
		Theme:	fakeTheme,
		ReceivedOn: time.Now().Unix(),
	}
	err = dispatcher.DispatchMessage(config.NewsiteQueueName, msg)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the Message to be produced
	time.Sleep(time.Second)
	expectedSitePath := filepath.Join(config.SitesDir, strconv.Itoa(int(fakeSiteId)))
	expectedSiteSourcePath := filepath.Join(expectedSitePath, "source")
	expectedSitePublicPath := filepath.Join(expectedSitePath, "public")

	// Test if source and public directories are created
	if _, err := os.Stat(expectedSiteSourcePath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(expectedSitePublicPath); err != nil {
		t.Fatal(err)
	}

	// Remove public directory to test re-generate
	if err := os.RemoveAll(expectedSitePublicPath); err != nil {
		t.Fatal(err)
	}

	// Dispatch a ReGenerate Message
	err = dispatcher.DispatchMessage(config.RegenerateQueueName, msg)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the Message to be produced
	time.Sleep(time.Second)

	// Test if public directories are re-created
	if _, err := os.Stat(expectedSitePublicPath); err != nil {
		t.Fatal(err)
	}
}

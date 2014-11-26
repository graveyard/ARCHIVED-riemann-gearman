package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Clever/gearadmin"
	"github.com/Clever/riemanner/riemanner"
	"github.com/amir/raidman"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

// DebugRaidmanClient sends events to a io.Writer instead of an actual Riemann server.
type DebugRaidmanClient struct {
	encoder *json.Encoder
}

// NewDebugRaidmanClient creates a DebugRaidmanClient.
func NewDebugRaidmanClient(output io.Writer) DebugRaidmanClient {
	return DebugRaidmanClient{encoder: json.NewEncoder(output)}
}

// Send processes a riemann event.
func (d DebugRaidmanClient) Send(event *raidman.Event) error {
	return d.encoder.Encode(event)
}

// writeStatus translates gearman status to a Riemann event.
func writeStatus(c riemanner.RaidmanClient, status gearadmin.Status, tags []string) error {
	if err := c.Send(&raidman.Event{
		Ttl:     10,
		Tags:    tags,
		Service: fmt.Sprintf("gearmand.status.%s.total", status.Function),
		Metric:  status.Total,
	}); err != nil {
		return err
	}
	if err := c.Send(&raidman.Event{
		Ttl:     10,
		Tags:    tags,
		Service: fmt.Sprintf("gearmand.status.%s.running", status.Function),
		Metric:  status.Running,
	}); err != nil {
		return err
	}
	return c.Send(&raidman.Event{
		Ttl:     10,
		Tags:    tags,
		Service: fmt.Sprintf("gearmand.status.%s.available", status.Function),
		Metric:  status.AvailableWorkers,
	})
}

func sendMetrics(g gearadmin.GearmanAdmin, r riemanner.RaidmanClient, tags []string) {
	statuses, err := g.Status()
	if err != nil {
		log.Fatalf("error retrieving gearman status: %s", err.Error())
	}
	for _, status := range statuses {
		if err := writeStatus(r, status, tags); err != nil {
			log.Fatalf("error writing status event: %s", err.Error())
		}
	}
}

func main() {
	var gearman, riemann, tagswithdelim string
	var interval int
	flag.StringVar(&gearman, "gearman", "tcp://localhost:4730",
		"Use the specified host:port to connect to gearman.")
	flag.StringVar(&riemann, "riemann", "tcp://localhost:5555",
		"Write events to Riemann running at this port. Can also specify 'stdout' to debug.")
	flag.StringVar(&tagswithdelim, "tags", "",
		"Tags to add to the Riemann event.")
	flag.IntVar(&interval, "interval", 60000, "Interval in ms to output data.")
	flag.Parse()
	tags := strings.Split(tagswithdelim, ",")

	// Set up gearadmin
	gearmanURL, err := url.Parse(gearman)
	if err != nil {
		log.Fatal("error parsing gearman url: %s", err)
	}
	gearmanConn, err := net.Dial(gearmanURL.Scheme, gearmanURL.Host)
	if err != nil {
		log.Fatal("error connecting to gearman: %s", err)
	}
	defer gearmanConn.Close()
	gearadmin := gearadmin.NewGearmanAdmin(gearmanConn)

	// Set up where we'll write riemann events
	var riemannClient riemanner.RaidmanClient
	if riemann == "stdout" {
		riemannClient = NewDebugRaidmanClient(os.Stdout)
	} else {
		riemannURL, err := url.Parse(riemann)
		if err != nil {
			log.Fatal("error parsing riemann url: %s", err)
		}
		raidmanClient, err := raidman.Dial(riemannURL.Scheme, riemannURL.Host)
		if err != nil {
			log.Fatal("error connecting to riemann:", err)
		}
		defer raidmanClient.Close()
		riemannClient = raidmanClient
	}

	// Send the stats when we first start, and then at the specified interval
	sendMetrics(gearadmin, riemannClient, tags)
	for _ = range time.Tick(time.Duration(interval) * time.Millisecond) {
		sendMetrics(gearadmin, riemannClient, tags)
	}
}

package main

import (
	"flag"
	"github.com/dhiaayachi/bambu-go"
	"github.com/dhiaayachi/bambu-go/events"
	"log"
	"time"
)

func main() {
	token := flag.String("token", "", "token")
	flag.Parse()
	if token == nil || *token == "" {
		flag.Usage()
		log.Fatalf("missing token")
	}
	b, err := bambu.NewBambuClient("us.mqtt.bambulab.com", "8883", *token, "https://api.bambulab.com")
	if err != nil {
		log.Fatal(err)
	}

	err = b.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = b.SubscribeAll(func(devID string, evt events.ReportEvent) {
		log.Printf("Received report event for dev_id=%s event_type=%s, event=%s", devID, evt.GetType(), evt.String())

	})
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(60 * time.Second)
}

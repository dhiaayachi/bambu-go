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

	doneCh := make(chan any)
	err = b.SubscribeAll(func(devID string, evt events.ReportEvent) {

		switch evt.GetType() {
		case events.SystemType:
			log.Printf("Received report event for dev_id=%s event_type=%s,event=%s", devID, evt.GetType(), evt.String())
			close(doneCh)
		default:
			log.Printf("Received report event for dev_id=%s event_type=%s", devID, evt.GetType())
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = b.UnsubscribeAll() }()

	devs := b.Devices()
	if len(devs) == 0 {
		log.Fatal("no devices found")
	}
	onoff := []string{"on", "off"}
	for i := 0; i < 10; i++ {
		err = b.Publish(devs[0], &events.LedCtlRequest{System: events.System{SequenceId: "99999", LedMode: onoff[i%2], Command: "ledctrl", LedNode: "chamber_light"}})
		if err != nil {
			log.Fatal(err)
		}
		<-doneCh
		time.Sleep(5 * time.Second)
		doneCh = make(chan any)
	}

}

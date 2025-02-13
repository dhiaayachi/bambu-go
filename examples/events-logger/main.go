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
	lan := flag.Bool("lan", false, "lan")
	host := flag.String("host", "", "host")
	port := flag.String("port", "8883", "host")
	user := flag.String("user", "bblp", "host")
	cert := flag.String("cert-file", "./ca_cert.pem", "host")
	flag.Parse()
	if token == nil || *token == "" {
		flag.Usage()
		log.Fatalf("missing token")
	}
	if *lan {
		if host == nil || *host == "" {
			flag.Usage()
			log.Fatalf("missing host")
		}
	}

	var err error
	var b *bambu.Client
	if *lan {
		b, err = bambu.NewBambuClientLan(*host, *port, *user, *token, *cert)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b, err = bambu.NewBambuClientCloud("us.mqtt.bambulab.com", "8883", *token, "https://api.bambulab.com")
		if err != nil {
			log.Fatal(err)
		}
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
	defer func() { _ = b.UnsubscribeAll() }()
	time.Sleep(60 * time.Second)
}

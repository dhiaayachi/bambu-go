package main

import (
	"flag"
	"github.com/dhiaayachi/bambu-go"
	"github.com/dhiaayachi/bambu-go/events"
	"log"
	"time"
)

func main() {
	token := flag.String("token", "", "token: the token to use for authentication in cloud mode or the access token in lan mode")
	lan := flag.Bool("lan", false, "lan mode (cloud mode if false)")
	host := flag.String("host", "", "host or IP address to connect to (only needed in lan mode)")
	port := flag.String("port", "8883", "port number")
	user := flag.String("user", "bblp", "username to connect to (only needed in lan mode)")
	cert := flag.String("cert-file", "./ca_cert.pem", "cert file to use to connect (only needed in lan mode)")
	serial := flag.String("serial", "", "printer serial number (only needed in lan mode)")
	flag.Parse()
	if token == nil || *token == "" {
		flag.Usage()
		log.Fatalf("missing token")
	}
	if port == nil || *port == "" {
		flag.Usage()
		log.Fatalf("missing host")
	}
	if *lan {
		if host == nil || *host == "" {
			flag.Usage()
			log.Fatalf("missing host")
		}
		if user == nil || *user == "" {
			flag.Usage()
			log.Fatalf("missing host")
		}
		if cert == nil || *cert == "" {
			flag.Usage()
			log.Fatalf("missing host")
		}
		if serial == nil || *serial == "" {
			flag.Usage()
			log.Fatalf("missing host")
		}
	}

	var err error
	var b *bambu.Client
	if *lan {
		b, err = bambu.NewBambuClientLan(*host, *port, *user, *token, *cert, *serial)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b, err = bambu.NewBambuClientCloud("us.mqtt.bambulab.com", *port, *token, "https://api.bambulab.com")
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

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

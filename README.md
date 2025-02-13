# Bambu Go Client

## Description

Bambu Go Client is a Go library designed to interact with BambuLab 3D printers via MQTT and HTTP. It allows subscribing to device reports, publishing events, and managing device connections seamlessly.

## Installation

To use this library in your Go project, you can install it using:

```sh
go get github.com/dhiaayachi/bambu-go
```

## Usage

### Get an account token

Use the provided login-helper tool to retrieve a valid token as follow

```shell
cd login-helper
go build .
./login-helper --account '<email>' --password '<makers world password>'
```
provide the access code, sent to your configured email for 2FA auth.

The returned access token will be similar to `Access Token: AACjxHoRLGcgKLlS8fQ....`  and is valid for 3 months.
pass that token to the bambu-go library.

Unfortunately, because of the 2FA limitation and the impossibility to get an access token without it, it's impossible to 
integrate the login flow in the library for now.

### Creating a Client (Cloud)

```go
package main

import (
	"fmt"
	"github.com/dhiaayachi/bambu-go/bambu"
)

func main() {
	client, err := bambu.NewBambuClientCloud("localhost", "1883", "your_token", "http://api.example.com")
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	fmt.Println("Client created successfully")
}
```

### Creating a Client (Cloud)

```go
package main

import (
	"fmt"
	"log"
	"github.com/dhiaayachi/bambu-go/bambu"
)

func main() {
	client, err := bambu.NewBambuClientLan("192.168.2.2", "8883", "bblp", "device_token", "./ca_cert.pem")
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	fmt.Println("Client created successfully")
}
```

### Subscribing to Device Events

```go
client.SubscribeAll(func(devID string, evt events.ReportEvent) {
	fmt.Printf("Received event from device %s: %v\n", devID, evt)
})
```

## Running Tests

To run the tests for this library, use:

```sh
go test ./...
```

## License

This project is licensed under the MIT License.


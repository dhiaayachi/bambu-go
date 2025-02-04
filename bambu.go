package bambu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dhiaayachi/bambu-go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/oklog/ulid/v2"
	"io"
	"net/http"
	"regexp"
	"time"
)

const preferenceUri = "/v1/design-user-service/my/preference"
const bindUri = "/v1/iot-service/api/user/bind"

type Device struct {
	DevId          string  `json:"dev_id"`
	Name           string  `json:"name"`
	Online         bool    `json:"online"`
	PrintStatus    string  `json:"print_status"`
	DevModelName   string  `json:"dev_model_name"`
	DevProductName string  `json:"dev_product_name"`
	DevAccessCode  string  `json:"dev_access_code"`
	NozzleDiameter float64 `json:"nozzle_diameter"`
	DevStructure   string  `json:"dev_structure"`
}

type devicesRsp struct {
	Message string      `json:"message"`
	Code    interface{} `json:"code"`
	Error   interface{} `json:"error"`
	Devices []Device    `json:"devices"`
}

type mqttClient interface {
	Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token
	Unsubscribe(topics ...string) mqtt.Token
	Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token
	Connect() mqtt.Token
}

type Client struct {
	mqttClient mqttClient
	host       string
	port       string
	username   string
	token      string
	apiUrl     string
	print      map[string][]byte
	devID      []string
}

func NewBambuClient(host string, port string, token string, url string) (*Client, error) {
	bambuClient := Client{host: host, port: port, apiUrl: url, token: token, print: make(map[string][]byte), devID: make([]string, 0)}

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", bambuClient.apiUrl, preferenceUri), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bambuClient.token))
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type u struct {
		Uid int `json:"uid"`
	}
	var uid u
	err = json.Unmarshal(body, &uid)
	if err != nil {
		return nil, err
	}

	bambuClient.username = fmt.Sprintf("u_%d", uid.Uid)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%s", bambuClient.host, bambuClient.port))
	opts.SetClientID(ulid.Make().String())
	opts.SetUsername(bambuClient.username)
	opts.SetPassword(bambuClient.token)
	bambuClient.mqttClient = mqtt.NewClient(opts)
	return &bambuClient, nil
}
func (b *Client) Connect() error {
	token := b.mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (b *Client) SubscribeAll(handler func(devId string, evt events.ReportEvent)) error {
	devices, err := b.getAllDevices()
	if err != nil {
		return err
	}

	for _, device := range devices {
		topic := fmt.Sprintf("device/%s/report", device)
		token := b.mqttClient.Subscribe(topic, 0, b.handlerWrapper(handler))
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}

		time.Sleep(1 * time.Second)
		// Push a pushall command to request a full update the first time
		req := events.PushRequest{SequenceId: "99999", Command: "pushall", Version: 1, PushTarget: 1}
		token = b.mqttClient.Publish(topic, 0, false, req.String())
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
		b.devID = append(b.devID, device)
	}
	return nil
}

func (b *Client) handlerWrapper(handler func(devId string, evt events.ReportEvent)) func(client mqtt.Client, message mqtt.Message) {
	return func(client mqtt.Client, message mqtt.Message) {
		devId, err := parseDevice(message.Topic())
		if err != nil {
			return
		}
		evtType := make(map[string]json.RawMessage)
		err = json.Unmarshal(message.Payload(), &evtType)
		if err != nil {
			return
		}
		for k, v := range evtType {
			var evt events.ReportEvent
			var newJ []byte
			switch k {
			case events.PrintType:
				// Print type provide an incremental update
				// This is handled by merging the event with an event cache locally
				var ok bool
				oldJ, ok := b.print[message.Topic()]
				if ok {
					newJ, err = jsonpatch.MergePatch(v, oldJ)
					if err != nil {
						return
					}
				} else {
					newJ = v
					b.print[message.Topic()] = v
				}
			default:
				newJ = v
			}
			evt = events.NewReportEvent(k)
			err := json.Unmarshal(newJ, evt)
			if err != nil {
				return
			}
			handler(devId, evt)
		}

	}
}

func (b *Client) getAllDevices() ([]string, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", b.apiUrl, bindUri), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", b.token))
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var u devicesRsp
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	if len(u.Devices) == 0 {
		return nil, fmt.Errorf("no devices found")
	}
	devs := make([]string, 0)
	for _, dev := range u.Devices {
		devs = append(devs, dev.DevId)
	}
	return devs, nil
}

func (b *Client) UnsubscribeAll() error {
	subscriptions := make([]string, len(b.devID))
	for i, dev := range b.devID {
		subscriptions[i] = fmt.Sprintf("device/%s/report", dev)
	}
	token := b.mqttClient.Unsubscribe(subscriptions...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (b *Client) Unsubscribe(devID string) error {
	token := b.mqttClient.Unsubscribe(devID)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (b *Client) Devices() []string {
	return b.devID
}

func parseDevice(topic string) (string, error) {
	re := regexp.MustCompile("device/(.*)/report")
	match := re.FindStringSubmatch(topic)
	if len(match) < 2 {
		return "", errors.New("invalid device id")
	}
	return match[1], nil
}

func (b *Client) Publish(devId string, evt events.RequestEvent) error {
	token := b.mqttClient.Publish(fmt.Sprintf("device/%s/request", devId), 0, false, evt.String())
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

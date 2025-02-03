package bambu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dhiaayachi/bambu-go/events"
	"github.com/dhiaayachi/bambu-go/jsonpatch"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oklog/ulid/v2"
	"io"
	"net/http"
	"regexp"
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
	print      map[string]*events.PrintReport
}

func NewBambuClient(host string, port string, token string, url string) (*Client, error) {
	bambuClient := Client{host: host, port: port, apiUrl: url, token: token, print: make(map[string]*events.PrintReport)}

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
func (b *Client) SubscribeAll(handler func(dev_id string, evt events.ReportEvent)) error {
	devices, err := b.getAllDevices()
	if err != nil {
		return err
	}
	for _, device := range devices {

		topic := "device/%s/report"

		token := b.mqttClient.Subscribe(fmt.Sprintf(topic, device), 0, func(client mqtt.Client, message mqtt.Message) {
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
				switch k {
				case events.PrintType:
					var ok bool
					evt, ok = b.print[message.Topic()]
					if !ok {
						b.print[message.Topic()] = &events.PrintReport{}
						evt = b.print[message.Topic()]
					}
					err := jsonpatch.PatchValues(v, evt.(*events.PrintReport))
					if err != nil {
						return
					}
				default:
					evt = events.NewReportEvent(k)
					err := json.Unmarshal(v, evt)
					if err != nil {
						return
					}

				}

				handler(devId, evt)
			}

		})
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
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
	token := b.mqttClient.Unsubscribe("dev")
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (b *Client) Unsubscribe(devID string) error {
	token := b.mqttClient.Unsubscribe(fmt.Sprintf("device/%s/report", devID), devID)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func parseDevice(topic string) (string, error) {
	re := regexp.MustCompile("device/(.*)/report")
	match := re.FindStringSubmatch(topic)
	if len(match) < 2 {
		return "", errors.New("invalid device id")
	}
	return match[1], nil
}

func (b *Client) Publish(devId string, evt *events.PrintReport) error {
	token := b.mqttClient.Publish(fmt.Sprintf("device/%s/request", devId), 0, false, evt)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

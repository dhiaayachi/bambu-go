package bambu

import (
	"errors"
	"github.com/dhiaayachi/bambu-go/events"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMqttClient struct {
	mock.Mock
}

type DummyToken struct {
	mock.Mock
}

func (t *DummyToken) Wait() bool {
	return t.Called().Bool(0)
}

func (t *DummyToken) WaitTimeout(duration time.Duration) bool {
	return t.Called(duration).Bool(0)
}

func (t *DummyToken) Done() <-chan struct{} {
	return t.Called().Get(0).(<-chan struct{})
}

func (t *DummyToken) Error() error {
	return t.Called().Error(0)
}

func (m *MockMqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	args := m.Called(topic, qos, callback)
	return args.Get(0).(mqtt.Token)
}

func (m *MockMqttClient) Unsubscribe(topics ...string) mqtt.Token {
	args := m.Called(topics)
	return args.Get(0).(mqtt.Token)
}

func (m *MockMqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	args := m.Called(topic, qos, retained, payload)
	return args.Get(0).(mqtt.Token)
}

func (m *MockMqttClient) Connect() mqtt.Token {
	args := m.Called()
	return args.Get(0).(mqtt.Token)
}

func TestNewBambuClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/design-user-service/my/preference", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"uid": 12345}`))
	}))
	defer server.Close()

	client, err := NewBambuClient("localhost", "1883", "test-token", server.URL)
	assert.NoError(t, err)
	assert.Equal(t, "u_12345", client.username)
}

func TestNewBambuClient_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := NewBambuClient("localhost", "1883", "test-token", server.URL)
	assert.Error(t, err)
}

func TestParseDevice(t *testing.T) {
	devID, err := parseDevice("device/abc123/report")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", devID)
}

func TestParseDevice_Invalid(t *testing.T) {
	_, err := parseDevice("invalid/topic")
	assert.Error(t, err)
}

func TestGetAllDevices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/iot-service/api/user/bind", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"devices": [{"dev_id": "device1"}, {"dev_id": "device2"}]}`))
	}))
	defer server.Close()

	client := &Client{apiUrl: server.URL, token: "test-token"}
	devices, err := client.getAllDevices()
	assert.NoError(t, err)
	assert.Equal(t, []string{"device1", "device2"}, devices)
}

func TestGetAllDevices_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"devices": []}`))
	}))
	defer server.Close()

	client := &Client{apiUrl: server.URL, token: "test-token"}
	_, err := client.getAllDevices()
	assert.Error(t, err)
}

func TestGetAllDevices_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &Client{apiUrl: server.URL, token: "test-token"}
	_, err := client.getAllDevices()
	assert.Error(t, err)
}

func TestConnect_Success(t *testing.T) {
	mockMqtt := new(MockMqttClient)
	mockToken := new(mqtt.DummyToken)
	mockMqtt.On("Connect").Return(mockToken)

	client := &Client{mqttClient: mockMqtt}
	err := client.Connect()
	assert.NoError(t, err)
	mockMqtt.AssertExpectations(t)
}

func TestConnect_Failure(t *testing.T) {
	mockMqtt := new(MockMqttClient)
	var mockToken = new(DummyToken)
	mockToken.On("Error").Return(errors.New("connection failed"))
	mockToken.On("Wait").Return(true)
	mockMqtt.On("Connect").Return(mockToken)

	client := &Client{mqttClient: mockMqtt}
	err := client.Connect()
	assert.Error(t, err)
	mockMqtt.AssertExpectations(t)
}

func TestSubscribeAll(t *testing.T) {
	mockMqtt := new(MockMqttClient)
	mockToken := new(mqtt.DummyToken)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/iot-service/api/user/bind", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"devices": [{"dev_id": "device1"}, {"dev_id": "device2"}]}`))
	}))
	defer server.Close()

	mockMqtt.On("Subscribe", "device/device1/report", byte(0), mock.Anything).Return(mockToken)
	mockMqtt.On("Subscribe", "device/device2/report", byte(0), mock.Anything).Return(mockToken)
	mockMqtt.On("Publish", "device/device1/request", byte(0), false, mock.Anything).Return(mockToken)
	mockMqtt.On("Publish", "device/device2/request", byte(0), false, mock.Anything).Return(mockToken)

	client := &Client{mqttClient: mockMqtt, apiUrl: server.URL, token: "test-token"}
	handler := func(devId string, evt events.ReportEvent) {}

	err := client.SubscribeAll(handler)
	assert.NoError(t, err)

	mockMqtt.AssertExpectations(t)
}

func TestHandlerWrapper(t *testing.T) {
	client := &Client{print: make(map[string][]byte)}
	handlerCalled := false

	var e events.ReportEvent
	handler := func(devId string, evt events.ReportEvent) {
		handlerCalled = true
		e = evt
	}

	wrappedHandler := client.handlerWrapper(handler)

	message := MessageMock{}
	message.On("Topic").Return("device/device1/report")
	message.On("Payload").Return([]byte(`{"print": {"command": "printing"}}`))

	wrappedHandler(nil, &message)

	assert.True(t, handlerCalled, "Handler should be called")

	assert.IsType(t, &events.PrintReport{}, e)
	pr := e.(*events.PrintReport)
	assert.Equal(t, "printing", pr.Command)

	message2 := MessageMock{}
	message2.On("Topic").Return("device/device1/report")

	message2.On("Payload").Return([]byte(`{"print": {"ams_status": 1}}`))

	wrappedHandler(nil, &message2)
	assert.IsType(t, &events.PrintReport{}, e)
	pr = e.(*events.PrintReport)
	assert.Equal(t, "printing", pr.Command)
	assert.Equal(t, 1, pr.AmsStatus)

}

type MessageMock struct {
	mock.Mock
}

func (m *MessageMock) Duplicate() bool {
	args := m.Mock.Called()
	return args.Bool(0)
}

func (m *MessageMock) Qos() byte {
	args := m.Mock.Called()
	return args.Get(0).(byte)
}

func (m *MessageMock) Retained() bool {
	args := m.Mock.Called()
	return args.Bool(0)
}

func (m *MessageMock) Topic() string {
	args := m.Mock.Called()
	return args.Get(0).(string)
}

func (m *MessageMock) MessageID() uint16 {
	args := m.Mock.Called()
	return uint16(args.Int(0))
}

func (m *MessageMock) Payload() []byte {
	args := m.Mock.Called()
	return args.Get(0).([]byte)
}

func (m *MessageMock) Ack() {
	m.Mock.Called()
}

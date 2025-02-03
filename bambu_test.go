package bambu

import (
	"github.com/dhiaayachi/bambu-go/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockMqttClient struct {
	mock.Mock
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

// Mock HTTP server response
func mockHTTPServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))
}

func TestNewClient(t *testing.T) {
	server := mockHTTPServer(`{"uid": "12345"}`, http.StatusOK)
	defer server.Close()

	client, err := NewBambuClient("localhost", "1883", "mockToken", server.URL)

	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "u_12345", client.username)
}

func TestNewClient_InvalidResponse(t *testing.T) {
	server := mockHTTPServer(`{"invalid": "response"}`, http.StatusOK)
	defer server.Close()

	client, err := NewBambuClient("localhost", "1883", "mockToken", server.URL)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestSubscribeAll(t *testing.T) {
	mockMqtt := new(MockMqttClient)
	bambu := &Client{
		mqttClient: mockMqtt,
		host:       "localhost",
		port:       "1883",
		username:   "user",
		token:      "mockToken",
		apiUrl:     "http://localhost",
	}

	mockMqtt.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(&mqtt.DummyToken{})
	testServer := mockHTTPServer(`{"devices": [{"dev_id": "dev1"}, {"dev_id": "dev2"}]}`, http.StatusOK)
	defer testServer.Close()

	bambu.apiUrl = testServer.URL
	err := bambu.SubscribeAll(func(dev_id string, evt events.ReportEvent) {})
	assert.NoError(t, err)
	mockMqtt.AssertExpectations(t)
}

func TestSubscribeAll_ErrorResponse(t *testing.T) {
	testServer := mockHTTPServer(`{"error": "invalid request"}`, http.StatusBadRequest)
	defer testServer.Close()

	bambu := &Client{
		host:   "localhost",
		port:   "1883",
		token:  "mockToken",
		apiUrl: testServer.URL,
	}

	err := bambu.SubscribeAll(func(dev_id string, evt events.ReportEvent) {})
	assert.Error(t, err)
}

func TestUnsubscribeAll(t *testing.T) {
	mockMqtt := new(MockMqttClient)
	mockToken := new(mqtt.DummyToken)
	mockMqtt.On("Unsubscribe", mock.Anything).Return(mockToken)

	bambu := &Client{mqttClient: mockMqtt}
	err := bambu.UnsubscribeAll()
	assert.NoError(t, err)
	mockMqtt.AssertExpectations(t)
}

func TestUnsubscribe(t *testing.T) {
	mockMqtt := new(MockMqttClient)
	mockToken := new(mqtt.DummyToken)
	mockMqtt.On("Unsubscribe", mock.Anything).Return(mockToken)

	bambu := &Client{mqttClient: mockMqtt}
	err := bambu.Unsubscribe("dev1")
	assert.NoError(t, err)
	mockMqtt.AssertExpectations(t)
}

func TestParseDevice(t *testing.T) {
	devId, err := parseDevice("device/abc123/report")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", devId)

	_, err = parseDevice("invalid/topic")
	assert.Error(t, err)
}

func TestPublish(t *testing.T) {
	mockMqtt := new(MockMqttClient)
	mockToken := new(mqtt.DummyToken)
	mockMqtt.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockToken)

	bambu := &Client{mqttClient: mockMqtt}
	evt := &events.PrintReport{}
	err := bambu.Publish("dev1", evt)
	assert.NoError(t, err)
	mockMqtt.AssertExpectations(t)
}

package mqttclient

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type topicLog struct {
	field string
	msg   string
	error error
}

func (l topicLog) logInfo() {
	msg := fmt.Sprintf("MQTT: %s - %s", l.field, l.msg)
	log.Info(msg)
}

func (l topicLog) logErr() {
	msg := fmt.Sprintf("MQTT: %s - %s - %v", l.field, l.msg, l.error)
	log.Info(msg)
}

// QOS describes the quality of service of an mqttClient publish
type QOS byte

const (
	// AtMostOnce means the broker will deliver at most once
	AtMostOnce QOS = iota
	// AtLeastOnce means the broker will deliver c message at least once
	AtLeastOnce
	// ExactlyOnce means the broker will deliver c message exactly once
	ExactlyOnce
)

// Client runs an mqttClient client
type Client struct {
	client                  mqtt.Client
	clientID                string
	connected               bool
	terminated              bool
	consumers               []consumer
	mqttPublishBuffers      []*MqttPublishBuffer
	mqttPublishBuffersMutex sync.Mutex
	chuckSize               int
}

// ClientOptions is the list of options used to create c client
type ClientOptions struct {
	Servers       []string // The list of broker hostnames to connect to
	ClientID      string   // If left empty c uuid will automatically be generated
	Username      string   // If not set then authentication will not be used
	Password      string   // Will only be used if the username is set
	AutoReconnect *bool    // If the client should automatically try to reconnect when the connection is lost
	ConnectRetry  *bool    // Automatically retry the connection in the event of a failure
}

type MqttPublishBuffer struct {
	Topic   string `json:"string"`
	Qos     QOS    `json:"qos"`
	Retain  bool   `json:"retain"`
	Payload string `json:"payload"`
}

type consumer struct {
	topic   string
	handler mqtt.MessageHandler
}

// Close agent
func (c *Client) Close() {
	c.client.Disconnect(250)
	c.terminated = true
}

// Subscribe to topic
func (c *Client) Subscribe(topic string, qos QOS, handler mqtt.MessageHandler) (err error) {
	token := c.client.Subscribe(topic, byte(qos), handler)
	if token.WaitTimeout(2*time.Second) == false {
		return errors.New("subscribe timout")
	}
	if token.Error() != nil {
		return token.Error()
	}
	c.consumers = append(c.consumers, consumer{topic, handler})
	return nil
}

// SubscribeMultiple subscribe to multiple topic
func (c *Client) SubscribeMultiple(filters map[string]byte, handler mqtt.MessageHandler) (err error) {
	token := c.client.SubscribeMultiple(filters, handler)
	if token.WaitTimeout(2*time.Second) == false {
		return errors.New("subscribe timout")
	}
	if token.Error() != nil {
		return token.Error()
	}
	for topic, _ := range filters {
		c.consumers = append(c.consumers, consumer{topic, handler})
	}
	return nil
}

// Unsubscribe unsubscribes from a certain topic and errors if this fails.
func (c *Client) Unsubscribe(topic string) error {
	token := c.client.Unsubscribe(topic)
	if token.Error() != nil {
		return token.Error()
	}
	return token.Error()
}

// Publish method buffers the list here and schedules it for later
func (c *Client) Publish(topic string, qos QOS, retain bool, payload string) {
	c.bufferMqttPublish(&MqttPublishBuffer{Topic: topic, Qos: qos, Retain: retain, Payload: payload})
}

func (c *Client) FlushMqttPublishBuffers() {
	log.Trace("Flushing mqtt publish buffers...")
	if len(c.mqttPublishBuffers) == 0 {
		log.Trace("MQTT publish buffers empty")
		return
	}
	c.mqttPublishBuffersMutex.Lock()
	chuckMqttPublishBuffers := ChuckMqttPublishBuffer(c.mqttPublishBuffers, c.chuckSize)
	c.mqttPublishBuffers = nil
	c.mqttPublishBuffersMutex.Unlock()
	for _, chuckMqttPublishBuffer := range chuckMqttPublishBuffers {
		wg := &sync.WaitGroup{}
		for _, record := range chuckMqttPublishBuffer {
			wg.Add(1)
			go func(record *MqttPublishBuffer) {
				defer wg.Done()
				log.Tracef("Publishing topic: %s", record.Topic)
				token := c.client.Publish(record.Topic, byte(record.Qos), record.Retain, record.Payload)
				if token.Error() != nil {
					log.Errorf("MQTT issue on publishing, topic: %s, error: %s", record.Topic, token.Error())
				}
			}(record)
			time.Sleep(1 * time.Millisecond) // for don't let them call at once
		}
		wg.Wait()
	}
	log.Trace("Finished MQTT publish buffers process")
}

// NewClient creates an mqttClient client
func NewClient(options ClientOptions, onConnected interface{}) (c *Client, err error) {
	c = &Client{chuckSize: 50}

	if options.ClientID == "" {
		options.ClientID, _ = nuuid.MakeUUID()
	}

	c.clientID = options.ClientID

	opts := mqtt.NewClientOptions()
	// brokers
	if options.Servers != nil && len(options.Servers) > 0 {
		for _, server := range options.Servers {
			opts.AddBroker(server)
		}
	} else {
		topicLog{"error", "min one server is required", nil}.logErr()
		return nil, err
	}

	if options.Username != "" {
		opts.SetUsername(options.Username)
		opts.SetPassword(options.Password)
	}

	opts.SetAutoReconnect(boolean.IsTrue(options.AutoReconnect))
	opts.SetConnectRetry(boolean.IsTrue(options.ConnectRetry))

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		topicLog{"error", "Lost connection", err}.logErr()
	}

	opts.OnConnect = func(cc mqtt.Client) {
		topicLog{"msg", "connected", nil}.logInfo()
		c.connected = true
		// Subscribe here, otherwise after connection lost,
		// you may not receive any message
		for _, s := range c.consumers {
			if token := cc.Subscribe(s.topic, 2, s.handler); token.Wait() && token.Error() != nil {
				topicLog{"error", "failed to subscribe", token.Error()}.logErr()
			}
			topicLog{"topic", "Resubscribe", nil}.logInfo()
		}
		if onConnected != nil {
			switch onConnected.(type) {
			case func():
				onConnected.(func())()
			}
		}
	}

	c.client = mqtt.NewClient(opts)

	go func() {
		done := make(chan os.Signal)
		<-done
		topicLog{"msg", "close down client", nil}.logInfo()
		c.Close()
	}()

	return c, nil
}

// Connect opens c new connection
func (c *Client) Connect() (err error) {
	token := c.client.Connect()
	if token.WaitTimeout(2*time.Second) == false {
		return errors.New("MQTT connection timeout")
	}
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}

func (c *Client) IsTerminated() bool {
	return c.terminated
}

func (c *Client) bufferMqttPublish(buffer *MqttPublishBuffer) {
	c.mqttPublishBuffersMutex.Lock()
	defer c.mqttPublishBuffersMutex.Unlock()
	for index, mpb := range c.mqttPublishBuffers {
		if mpb.Topic == buffer.Topic {
			c.mqttPublishBuffers[index] = buffer
			return
		}
	}
	c.mqttPublishBuffers = append(c.mqttPublishBuffers, buffer)
}

func (c *Client) removeIndex(index int) []*MqttPublishBuffer {
	return append(c.mqttPublishBuffers[:index], c.mqttPublishBuffers[index+1:]...)
}

func ChuckMqttPublishBuffer(array []*MqttPublishBuffer, chunkSize int) [][]*MqttPublishBuffer {
	var chucks [][]*MqttPublishBuffer
	for i := 0; i < len(array); i += chunkSize {
		end := i + chunkSize
		if end > len(array) {
			end = len(array)
		}
		chucks = append(chucks, array[i:end])
	}
	return chucks
}

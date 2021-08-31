package mqttClient

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)


type topicLog struct {
	field  string
	msg string
	error error
}

func (l topicLog) printLog() {
	msg := fmt.Sprintf("log: %s - %s -%v", l.field, l.msg, l.error)
	log.Println(msg)
}

// QOS describes the quality of service of an mqttClient publish
type QOS byte

const (
	// AtMostOnce means the broker will deliver at most once to every producer - this means message delivery is not guaranteed
	AtMostOnce QOS = iota
	// AtLeastOnce means the broker will deliver c message at least once to every producer
	AtLeastOnce
	// ExactlyOnce means the broker will deliver c message exactly once to every producer
	ExactlyOnce
)


// Client runs an mqttClient client
type Client struct {
	client        	mqtt.Client
	clientID      	string
	connected    	bool
	terminated    	bool
	consumers 	[]consumer
}


// ClientOptions is the list of options used to create c client
type ClientOptions struct {
	Servers  []string // The list of broker hostnames to connect to
	ClientID string   // If left empty c uuid will automatically be generated
	Username string   // If not set then authentication will not be used
	Password string   // Will only be used if the username is set
	SetKeepAlive time.Duration
	SetPingTimeout  time.Duration
	AutoReconnect bool // If the client should automatically try to reconnect when the connection is lost
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

// Unsubscribe unsubscribes from a certain topic and errors if this fails.
func (c *Client) Unsubscribe(topic string) error {
	token := c.client.Unsubscribe(topic)
	if token.Error() != nil {
		return token.Error()
	}
	return token.Error()
}

//// SubscribeMultiple subscribes to multiple topics and errors if this fails.
//func (c *Client) SubscribeMultiple(ctx context.Context, consumers map[string]QOS) error {
//	subs := make(map[string]byte, len(consumers))
//	for topic, qos := range consumers {
//		subs[topic] = byte(qos)
//	}
//	token := c.client.SubscribeMultiple(subs, nil)
//	err := tokenWithContext(ctx, token)
//	return err
//}

// Publish things
func (c *Client) Publish(topic string, qos QOS, retain bool, payload string) (err error) {
	token := c.client.Publish(topic, byte(qos), retain, payload)
	if token.WaitTimeout(2*time.Second) == false {
		return errors.New("publish timout")
	}
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

// NewClient creates an mqttClient client
func NewClient(options ClientOptions) (c *Client) {
	c = &Client{}
	opts := mqtt.NewClientOptions()
	// brokers
	if options.Servers != nil && len(options.Servers) > 0 {
		for _, server := range options.Servers {
			fmt.Println(server)
			opts.AddBroker(server)
		}
	} else {
		topicLog{"error", "min one server is required", nil}.printLog()
	}

	if options.ClientID == "" {
		options.ClientID, _ = utils.MakeUUID()
	}

	c.clientID = options.ClientID
	if options.Username != "" {
		opts.SetUsername(options.Username)
		opts.SetPassword(options.Password)
	}
	if options.SetKeepAlive == 0 {
		options.SetKeepAlive = 5
	}
	if options.SetPingTimeout == 0 {
		options.SetPingTimeout = 5
	}

	opts.SetAutoReconnect(options.AutoReconnect)
	opts.SetKeepAlive(options.SetKeepAlive * time.Second)
	opts.SetPingTimeout(options.SetPingTimeout * time.Second)

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		topicLog{"error", "Lost connection", nil}.printLog()
	}
	opts.OnConnect = func(cc mqtt.Client) {
		topicLog{"msg", "connected", nil}.printLog()
		c.connected = true
		//Subscribe here, otherwise after connection lost,
		//you may not receive any message
		for _, s := range c.consumers {
			if token := cc.Subscribe(s.topic, 2, s.handler); token.Wait() && token.Error() != nil {
				topicLog{"error", "failed to subscribe", token.Error()}.printLog()
			}
			topicLog{"topic", "Resubscribe", nil}.printLog()
		}
	}
	c.client = mqtt.NewClient(opts)
	go func() {
		done := make(chan os.Signal)
		<-done
		topicLog{"msg", "close down client", nil}.printLog()
		c.Close()
	}()

	return c
}

// Connect opens c new connection
func (c *Client) Connect() (err error) {
	token := c.client.Connect()
	if token.WaitTimeout(2 * time.Second) == false {
		return errors.New("open timeout")
	}
	if token.Error() != nil {
		return token.Error()
	}

	return
}

func (c *Client) IsConnected() bool {
	return c.connected
}

func (c *Client) IsTerminated() bool {
	return c.terminated
}
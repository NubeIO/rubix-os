package eventbus

import (
	"context"
	"fmt"
	"github.com/mustafaturan/bus/v3"
	"log"
)

// BusService :nodoc:
type BusService interface {
	EmitString(ctx context.Context, topicName string, data string) error
	Emit(ctx context.Context, topicName string, data interface{}) error
	RegisterTopic(topic string)
	RegisterTopicParent(parent string, child string)
}

var bu *bus.Bus
var busContext context.Context
var bs BusService

func Init() {
	bu = newBus()
	busContext = context.Background()
	bu.RegisterTopics(BusTopics()...)
}

func GetBus() *bus.Bus {
	return bu
}

func CTX() context.Context {
	return busContext
}

func GetService() BusService {
	return bs
}

type notificationService struct {
	eb *bus.Bus
}

// NewService ...
func NewService(eb *bus.Bus) BusService {
	ns := &notificationService{
		eb: eb,
	}
	bs = ns
	return ns
}

// EmitString emits an event to the bus
func (eb *notificationService) EmitString(ctx context.Context, topicName string, data string) error {
	ctx = context.WithValue(ctx, bus.CtxKeyTxID, "")
	err := eb.eb.Emit(ctx, topicName, data)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return err
}

// Emit emits an event to the bus
func (eb *notificationService) Emit(ctx context.Context, topicName string, data interface{}) error {
	err := eb.eb.Emit(ctx, topicName, data)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return err
}

// RegisterTopic registers a topic for publishing
func (eb *notificationService) RegisterTopic(ds string) {
	eb.eb.RegisterTopics(fmt.Sprintf("%s", ds))
}

// RegisterTopicParent registers a topic for publishing
func (eb *notificationService) RegisterTopicParent(parent string, child string) {
	topic := fmt.Sprintf("%s.%s", parent, child)
	eb.eb.RegisterTopics(topic)
}

// UnregisterTopic removes a topic from consumer
func (eb *notificationService) UnregisterTopic(topic string) {
	eb.eb.DeregisterTopics(topic)
}

// UnregisterTopicChild removes a topic from consumer
func (eb *notificationService) UnregisterTopicChild(parent string, child string) {
	topic := fmt.Sprintf("%s.%s", parent, child)
	eb.eb.DeregisterTopics(topic)
}

// UnsubscribeHandler removes handler
func (eb *notificationService) UnsubscribeHandler(id string) {
	eb.eb.DeregisterHandler(id)
}

package eventbus

import (
	"context"
	"fmt"
	"github.com/mustafaturan/bus/v3"
)

type notificationService struct {
	eb   *bus.Bus
}


// NewNotificationService ...
func NewNotificationService(eb *bus.Bus) NotificationService {
	ns := &notificationService{
		eb: eb,
	}
	ns.registerPointsSubscriber() //add as types needed
	return ns
}


// EmitString emits an event to the bus
func (eb *notificationService) EmitString(ctx context.Context, topicName string, data string) {
	ctx = context.WithValue(ctx, bus.CtxKeyTxID, "")
	err := eb.eb.Emit(ctx, topicName, data)
	fmt.Println("Emit")
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Emit emits an event to the bus
func (eb *notificationService) Emit(ctx context.Context, topicName string, data interface{}) {
	err := eb.eb.Emit(ctx, topicName, data)
	fmt.Println("Emit")
	if err != nil {
		fmt.Println(err.Error())
	}
}



func AddTopic(eb *bus.Bus) NotificationService {
	ns := &notificationService{
		eb: eb,
	}
	ns.registerPointsSubscriber() //add as types needed
	return ns
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

// UnregisterTopic removes a topic from subscription
func (eb *notificationService) UnregisterTopic(topic string) {
	eb.eb.DeregisterTopics(topic)
}


// UnregisterTopicChild removes a topic from subscription
func (eb *notificationService) UnregisterTopicChild(parent string, child string) {
	topic := fmt.Sprintf("%s.%s", parent, child)
	eb.eb.DeregisterTopics(topic)
}

// UnsubscribeHandler removes handler
func (eb *notificationService) UnsubscribeHandler(id string) {
	eb.eb.DeregisterHandler(id)
}



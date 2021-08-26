package eventbus

import (
	"context"
	"fmt"
	"github.com/mustafaturan/bus/v3"
	"github.com/mustafaturan/monoton/v2"
	"github.com/mustafaturan/monoton/v2/sequencer"
)


type DatasetName struct {
	Name string `json:"Name"`
}

type EventBus interface {
	Init(datasets []DatasetName)
	RegisterTopic(ds string)
	UnregisterTopic(ds string)
	SubscribeToDataset(id string, matcher string, f func(ctx context.Context, e bus.Event))
	UnsubscribeToDataset(id string)
	Emit(ctx context.Context, topicName string, data interface{})
}

type MEventBus struct {
	Bus    *bus.Bus

}

func NewBus2() (EventBus, error) {
	// configure id generator
	node := uint64(1)
	initialTime := uint64(1577865600000) // set 2020-01-01 PST as initial time
	m, err := monoton.New(sequencer.NewMillisecond(), node, initialTime)
	if err != nil {
		return nil, err
	}
	// init an id generator
	var idGenerator bus.Next = m.Next
	// create a new bus instance
	b, err := bus.NewBus(idGenerator)
	if err != nil {
		return nil, err
	}

	eb := &MEventBus{
		Bus:    b,
	}
	return eb, nil
}


func (eb *MEventBus) Init(datasets []DatasetName) {
	for _, ds := range datasets {
		eb.Bus.RegisterTopics(fmt.Sprintf("dataset.%s", ds.Name))
	}
}

// RegisterTopic registers a topic for publishing. "dataset." is prefixed in front of the topic
func (eb *MEventBus) RegisterTopic(ds string) {
	eb.Bus.RegisterTopics(fmt.Sprintf("dataset.%s", ds))
}

// UnregisterTopic removes a topic from subscription
func (eb *MEventBus) UnregisterTopic(ds string) {
	topic := fmt.Sprintf("dataset.%s", ds)
	eb.Bus.DeregisterTopics(topic)
}

// SubscribeToDataset adds a subscription to an already registered topic.
// The id should be unique, the matcher is a regexp to match against the registered topics,
// ie: dataset.*, dataset.sdb.*, dataset.sdb.Animal are all valid registrations.
// f is the func to be called
func (eb *MEventBus) SubscribeToDataset(id string, matcher string, f func(ctx context.Context, e bus.Event)) {
	//eb.logger.Infof("Registering subscription '%s' with matcher 'dataset.%s'", id, matcher)
	handler := bus.Handler{
		Handle:  f,
		Matcher: "dataset." + matcher,
	}
	eb.Bus.RegisterHandler(id, handler)
}

// UnsubscribeToDataset removes a dataset subscription
func (eb *MEventBus) UnsubscribeToDataset(id string) {
	eb.Bus.DeregisterHandler(id)
}

// Emit emits an event to the bus
func (eb *MEventBus) Emit(ctx context.Context, topicName string, data interface{}) {
	ctx = context.WithValue(ctx, bus.CtxKeyTxID, "")
	err := eb.Bus.Emit(ctx, topicName, data)
	if err != nil {
		fmt.Println(err.Error())
	}
}

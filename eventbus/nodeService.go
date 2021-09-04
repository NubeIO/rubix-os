package eventbus

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/mustafaturan/bus/v3"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

var C = cache.New(5*time.Minute, 10*time.Minute)

func getUnderlyingAsValue(data interface{}) reflect.Value {
	return reflect.ValueOf(data)
}
func setOutTopic(ioNum string,uuid string) string {
	return fmt.Sprintf("in%s.%s", ioNum ,uuid)
}
func setInputValue(topic string, payload string)  {
	C.Set(topic, payload, cache.DefaultExpiration)
}
func getInputValue(topic string) interface{} {
	v, _ := C.Get(topic)
	return v
}

func (eb *notificationService) registerNodes() {
	handler := bus.Handler{
		Handle:  func(ctx context.Context, e bus.Event){
			go func() {
			switch e.Topic {
			case NodeEventIn: //from db event
				payload, ok := e.Data.(*model.NodeList);if !ok {
					logrus.Error("EVENTBUS NodeEventIn failed to pass in payload")
					return
				}
				var node  *model.NodeList
				if x, found := C.Get(payload.UUID); found {
					node = x.(*model.NodeList)
				}

				in1Updated := false
				in1Topic := setOutTopic("1", node.UUID)
				in2Updated := false
				in2Topic := setOutTopic("2", node.UUID)
				if node.In1 != "" {
					if getInputValue(in1Topic) != node.In1 {
						in1Updated = true
						setInputValue(in1Topic, node.In1)
					}
				}
				if node.In2 != "" {
					if getInputValue(in2Topic) != node.In2 {
						in2Updated = true
						setInputValue(in2Topic, node.In2)
					}
				}
				fmt.Println("!!!!!! NEW INPUT EVENT", node.UUID, "in1Topic", node.Name, "payload", node.In1, "in1Updated", in1Updated)
				if in1Updated || in2Updated {
						k1 := fmt.Sprintf("in1.%s", node.UUID)
						k2 := fmt.Sprintf("in2.%s", node.UUID)
						in1, _ := C.Get(k1)
						in2, _ := C.Get(k2)
						s1 := getUnderlyingAsValue(in1).String()
						s2 := getUnderlyingAsValue(in2).String()
						out := ""
						out1Topic := ""
						if node.NodeType == "add" {
							out = add(s1, s2) //TODO update its context
							for _, el := range node.NodeOut1 {
								var updateNode  *model.NodeList
								out1Topic = fmt.Sprintf("out1.%s", node.UUID) //set out1 topic
								if x, found := C.Get(el.ToUUID); found {
									updateNode = x.(*model.NodeList)
								}
								updateNode.In1 = out
								fmt.Println("!!!!!!!!! OUT", "TOPIC FROM", node.Name, "out1Topic", out1Topic, el.UUID)
								node.Out1Value = out
								eventOut(updateNode)
							}
						}
				}
			case NodeEventOut:
				payload, ok := e.Data.(*model.NodeList)
				msg := fmt.Sprintf("out event%s", payload.Name)
				logrus.Info(msg)
				if !ok {
					return
				}
			}
			}()
		},
		Matcher: NodesAll,
	}
	eb.eb.RegisterHandler("nodes.*", handler)
}

func eventOut(body *model.NodeList){
	n := NewNotificationService(BUS)
	n.Emit(BusContext, NodeEventIn, body)
}


func add(in1 string, in2 string) string {
	return in1+in2
}
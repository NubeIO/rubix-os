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

var c = cache.New(5*time.Minute, 10*time.Minute)

func getUnderlyingAsValue(data interface{}) reflect.Value {
	return reflect.ValueOf(data)
}
func setOutTopic(ioNum string,uuid string) string {
	return fmt.Sprintf("in%s.%s", ioNum ,uuid)
}
func setInputValue(topic string, payload string)  {
	c.Set(topic, payload, cache.DefaultExpiration)
}
func getInputValue(topic string) interface{} {
	v, _ := c.Get(topic)
	return v
}

func (eb *notificationService) registerNodes() {
	handler := bus.Handler{
		Handle:  func(ctx context.Context, e bus.Event){
			go func() {
			switch e.Topic {
			case NodeEventIn: //from db event
				fmt.Println("NodeEventIn")
				payload, ok := e.Data.(*model.NodeList)
				in1Updated := false
				in1Topic := setOutTopic("1", payload.UUID)
				in2Updated := false
				in2Topic := setOutTopic("2", payload.UUID)
				if payload.In1 != "" {
					if getInputValue(in1Topic) != payload.In1 {
						in1Updated = true
						setInputValue(in1Topic, payload.In1)
					}
				}
				if payload.In2 != "" {
					if getInputValue(in2Topic) != payload.In2 {
						in2Updated = true
						setInputValue(in2Topic, payload.In2)
					}
				}
				if in1Updated || in2Updated {
					if payload.In1 != "null" && payload.In2 != "null"{
						k1 := fmt.Sprintf("in1.%s", payload.UUID)
						k2 := fmt.Sprintf("in2.%s", payload.UUID)
						in1, _ := c.Get(k1)
						in2, _ := c.Get(k2)
						s1 := getUnderlyingAsValue(in1).String()
						s2 := getUnderlyingAsValue(in2).String()
						out1Updated := false
						out := ""
						in1NewValueTopic := ""
						out1Topic := ""
						if payload.NodeType == "add" {
							out = add(s1, s2)
							for _, el := range payload.NodeOut1 {
								out1Topic = fmt.Sprintf("out1.%s", payload.UUID) //set out1 topic
								o1, _ := c.Get(k1)
								if o1 != out {
									c.Set(out1Topic, out, cache.NoExpiration) // fire outputs
									in1NewValueTopic = fmt.Sprintf("%s.%s", el.Connection, el.ToUUID) //set input topic
									out1Updated = true
								}
								if out1Updated {
									eventOut(in1NewValueTopic, out)
									fmt.Println("RE-SYNC INPUTS new payload from UUID", out1Topic, "to", in1NewValueTopic, "value", out)
								}
							}
						}
					}
				}
				if !ok {
					return
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

func eventOut(UUID string, data string){
	n := NewNotificationService(BUS)
	d := new(model.NodeList)
	d.In1 = data
	d.UUID = UUID
	fmt.Println("!!!!  OUT EVENT !!!!")
	n.Emit(BusContext, NodeEventIn, d)
}


func add(in1 string, in2 string) string {
	return in1+in2
}
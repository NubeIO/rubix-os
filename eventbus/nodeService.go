package eventbus

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
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
				fmt.Println("INPUT EVENT FROM NAME:", node.Name, "UUID", node.UUID, "IN1", node.In1, "IN2", node.In2, "in1Updated", in1Updated, "in2Updated", in2Updated)
				if in1Updated || in2Updated {
						if node.NodeType == "add" {
							add(node)
						} else 	if node.NodeType == "addDog" {
							addDog(node)
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


//add node adds to string(its a demo node)
func add(node *model.NodeList){
	fmt.Println("ADD-NODE")
	k1 := fmt.Sprintf("in1.%s", node.UUID)
	k2 := fmt.Sprintf("in2.%s", node.UUID)
	in1, _ := C.Get(k1)
	in2, _ := C.Get(k2)
	s1 := getUnderlyingAsValue(in1).String()
	s2 := getUnderlyingAsValue(in2).String()
	out := ""
	out1Topic := "" //TODO make it update itself out1 (or all outputs that are used by the node developer)
	//if this node has 1 or more outputs to the same node then the calc must be done before send the output
	//loop through a node and get the out connections that are for the same node
	//then do the calc and send the output to the node
	var outList []interface{}
		for _, el := range node.NodeOut1 { //update the node context
			var updateNode  *model.NodeList
			out1Topic = fmt.Sprintf("out1.%s", node.UUID)
			if x, found := C.Get(el.ToUUID); found {
				updateNode = x.(*model.NodeList)
			}
			out = s1 + s2
			if el.Connection == "in1"{
				updateNode.In1 = out
			}
			if el.Connection == "in2"{
				updateNode.In2 = out
			}
			list := utils.NewArray()
			list.AddIfNotExist(el.ToUUID)
			outList = list.Values()
			node.Out1Value = out
		}

		for _, el := range outList { //publish the updated nodes on the bus
		var updateNode  *model.NodeList
		out1Topic = fmt.Sprintf("out1.%s", node.UUID)
		if x, found := C.Get(el.(string)); found {
			updateNode = x.(*model.NodeList)
		}
		eventOut(updateNode)
		fmt.Println("OUT ADD", "FROM-NODE", node.Name, "out1Topic", out1Topic, "TO-NODE", updateNode.Name, "TO In1", k1, "TO In2", k2, "OUT", out)
	}

}


//addDog add I like dogs to the incoming string
func addDog(node *model.NodeList){
	k1 := fmt.Sprintf("in1.%s", node.UUID)
	in1, _ := C.Get(k1)
	s1 := getUnderlyingAsValue(in1).String()
	out := ""
	out1Topic := "" //TODO make it update itself out1 (or all outputs that are used by the node developer)
	for _, el := range node.NodeOut1 {
		var updateNode  *model.NodeList
		out1Topic = fmt.Sprintf("out1.%s", node.UUID)
		if x, found := C.Get(el.ToUUID); found {
			updateNode = x.(*model.NodeList)
		}
		if el.Connection == "in1"{
			updateNode.In1 = out
		}
		if el.Connection == "in2"{
			updateNode.In2 = out
		}
		fmt.Println("!!!!!!!!! OUT", "TOPIC FROM", node.Name, "out1Topic", out1Topic, el.UUID, updateNode.In1, updateNode.In2)
		out = s1 + "I like dogs"
		node.Out1Value = out
		eventOut(updateNode)
	}
}
package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
)

func main()  {

	c := client.NewFlowRestClient("admin", "admin")


	//token, err := c.GetToken("admin", "admin")
	addNet, err := c.ClientAddNetwork()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Add network")
	fmt.Println(addNet.Status)
	fmt.Println(addNet.Response.UUID)
	fmt.Println(addNet.Response.Name)

	addDev, err := c.ClientAddDevice(addNet.Response.UUID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Add device")
	fmt.Println(addDev.Status)
	fmt.Println(addDev.Response.UUID)
	fmt.Println(addDev.Response.Name)

	addPoint, err := c.ClientAddPoint(addDev.Response.UUID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Add point")
	fmt.Println(addPoint.Status)
	fmt.Println(addPoint.Response.UUID)
	fmt.Println(addPoint.Response.Name)


	addPoint2, err := c.ClientAddPoint(addDev.Response.UUID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Add point 2")
	fmt.Println(addPoint2.Status)
	fmt.Println(addPoint2.Response.UUID)
	fmt.Println(addPoint2.Response.Name)


	addGateway, err := c.ClientAddGateway(false)
	if err != nil {
		fmt.Println(err)
		return

	}
	fmt.Println("Add gateway")
	fmt.Println(addGateway.Status)
	fmt.Println(addGateway.Response.UUID)
	fmt.Println(addGateway.Response.Name)

	// point 2 to make a subscriber connection to point 1
	tSub := new(client.Subscriber)
	tSub.Name = "test"
	tSub.Enable = true
	tSub.FromUUID = addPoint2.Response.UUID //from point 2
	tSub.ToUUID = addPoint.Response.UUID  //to point 1
	tSub.GatewayUuid = addGateway.Response.UUID
	tSub.SubscriberApplication = "mapping"
	tSub.SubscriberType = "point"

	addSubscriber, err := c.ClientAddSubscriber(*tSub)
	if err != nil {
		fmt.Println(err)
		return

	}
	fmt.Println("Add Subscriber")
	fmt.Println(addSubscriber.Status)
	fmt.Println(addSubscriber.Response.UUID)
	fmt.Println(addSubscriber.Response.Name)


	fmt.Println("FLOW-FRAMEWORK-TOKEN", c.ClientToken)


	//// point 1 to have 1 subscription to point 2
	//tSub2 := new(client.Subscription)
	//tSub2.Name = "test"
	//tSub2.Enable = true
	//tSub2.ThingUuid = addPoint.Response.UUID //pass in point 1 UUID
	//tSub2.GatewayUuid = addGateway.Response.UUID
	//tSub2.SubscriberApplication = "mapping"
	//tSub2.SubscriberType = "point"
	//
	//addSubscription, err := c.ClientAddSubscription(*tSub2)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//
	//}
	//fmt.Println("Add Subscription")
	//fmt.Println(addSubscription.Status)
	//fmt.Println(addSubscription.Response.UUID)
	//fmt.Println(addSubscription.Response.Name)








}


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


	addGateway, err := c.ClientAddGateway(true)
	if err != nil {
		fmt.Println(err)
		return

	}
	fmt.Println("Add gateway")
	fmt.Println(addGateway.Status)
	fmt.Println(addGateway.Response.UUID)
	fmt.Println(addGateway.Response.Name)

	tSub := new(client.Subscriber)
	tSub.Name = "test"
	tSub.Enable = true
	tSub.ThingUuid = addPoint.Response.UUID
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








}


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







}


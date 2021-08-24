package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
)

func main()  {

	c := client.NewFlowRestClient("admin", "admin", "0.0.0.0", "1660")

	getPlat, err := c.ClientGetRubixPlat()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("getPlat")
	fmt.Println(getPlat.Status)
	fmt.Println(getPlat.Response.UUID)
	fmt.Println(getPlat.Response.Name)

	fmt.Println("FLOW-FRAMEWORK-TOKEN", c.ClientToken)








}


package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)



func main()  {

	c := client.NewSessionWithToken("CVvjlk3SwM7p4hp", "0.0.0.0", "1660")
	//
	p := new(model.Producer)
	p.Name = "new 2222"
	uuid := "pro_772748e553684000"
	pnt, err := c.ClientEditProducer(uuid, *p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("edit point")
	fmt.Println(pnt)

	//pntGet, err := c.ClientGetPoint("pnt_be1c5e0c1bad43c9")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("edit point")
	//fmt.Println(pntGet.Points.Name)

}


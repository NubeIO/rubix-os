package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)



func main()  {

	c := client.NewSessionWithToken("CjtuhBWiHXPLgO1", "0.0.0.0", "1660")
	//
	p := new(model.WriterClone)
	p.WriteValue = 22
	uuid := "sul_3542ab1febc24ddf"
	//pnt, err := c.ClientGetWriterClone(uuid, *p)
	pnt, err := c.ClientEditWriterClone(uuid, *p)
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


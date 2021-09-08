package main

import (
	"fmt"
	"strings"
)

func getTopicPart(topic string, index int, contains string) string {
	s := strings.Split(topic, ".")
	for i, e := range s {
		if i == index {
			if strings.Contains(e, contains) { // if topic has pnt (is uuid of point)
				return e
			}
		}
	}
	return ""
}

func main() {
	fmt.Println(strings.Contains("dev_da08b647adca44e0", "dev"))
	//c := client.NewSession("admin", "admin", "0.0.0.0", "1660")
	//
	//getPlat, err := c.ClientGetRubixPlat()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("getPlat")
	//fmt.Println(getPlat.Status)
	//fmt.Println(getPlat.Response.UUID)
	//fmt.Println(getPlat.Response.Name)
	//
	//fmt.Println("FLOW-FRAMEWORK-TOKEN", c.ClientToken)

	s := strings.Split("plugin.updated.plg_caf8c499eda74a84.dev_da08b647adca44e0", ".")
	for i, e := range s {
		if i == 3 {
			fmt.Println(i, e)
		}

	}
	fmt.Println(s)

}

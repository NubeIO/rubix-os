package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/src/system/networking"
	unit "github.com/NubeDev/flow-framework/src/units"
	"github.com/NubeDev/flow-framework/utils"
	"reflect"
	"time"
)

func main() {

	_, res, err := unit.Process(1, "c", "c")
	if err != nil {
		return
	}
	fmt.Println(err)
	fmt.Println(res.String())
	fmt.Println(res.AsFloat())
	fmt.Println(utils.RandInt(1, 11))
	fmt.Println(utils.RandInt(1, 11))
	fmt.Println(utils.RandFloat(1, 1011))
	fmt.Println(unit.Exists("length1"))

	aa := client.NewSessionNoAUTH("0.0.0.0", 1660)
	ping, err := aa.Ping()
	if err != nil {
		return
	}
	file := "/tmp/test.json"
	i := reflect.ValueOf(ping).Interface().(interface{})
	_, err = utils.WriteDataToFileAsJSON(i, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	t := time.Now()
	fmt.Println(t.Format("2006-01-02-15-04-05"))

	aaa, _ := networking.GetValidNetInterfacesForWeb()
	fmt.Println(aaa)
	//fmt.Println(networking.ExternalIpAddress())

}

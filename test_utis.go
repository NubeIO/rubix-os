package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
)

func main(){

	tf := utils.ToFloat64("123")
	fmt.Println(tf)

	mArr := utils.NewArray()
	fmt.Println(mArr.RemoveNil(mArr.Add(1,2,1,nil).Values()))
	fmt.Println(mArr.Add(false, false).Exist(true))
	fmt.Println(mArr.Add(true, false).Exist(true))
	fmt.Println(utils.ToFloat64("22")+11)
	fmt.Println(utils.Round(123.123, 1))
	fmt.Println(utils.Round(123.523, 1))
	fmt.Println(utils.RoundTo(123.123, 1))
	//COVs
	fmt.Println(utils.COV(1.1, 1.1, 1))
	fmt.Println(utils.COV(1.1, 2.0, 1))
	fmt.Println(utils.COV(1.1, 2.1, 1))
	fmt.Println(utils.COV(0, 1.0, 1))
	fmt.Println(utils.COV(0.1, 1.0, 1))
}

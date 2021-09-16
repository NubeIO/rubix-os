package main

import (
	"github.com/NubeDev/flow-framework/utils"
	"strings"
)



func topicParts(topic string) *utils.Array {
	s := strings.SplitAfter(topic, "/")
	arr := utils.NewArray()
	for _, e := range s {
		res := strings.ReplaceAll(e, "/", "")
		if res != ""{
			arr.Add(res)
		}
	}
	return arr
}


func main(){

	//tf := utils.ToFloat64("123")
	//fmt.Println(tf)
	//
	//mArr := utils.NewArray()
	//fmt.Println(mArr.RemoveNil(mArr.Add(1,2,1,nil).Values()))
	//fmt.Println(mArr.Add(false, false).Exist(true))
	//fmt.Println(mArr.Add(true, false).Exist(true))
	//fmt.Println(utils.ToFloat64("22")+11)
	//fmt.Println(utils.Round(123.123, 1))
	//fmt.Println(utils.Round(123.523, 1))
	//fmt.Println(utils.RoundTo(123.123, 1))
	////COVs
	//fmt.Println(utils.COV(1.1, 1.1, 1))
	//fmt.Println(utils.COV(1.1, 2.0, 1))
	//fmt.Println(utils.COV(1.1, 2.1, 1))
	//fmt.Println(utils.COV(0, 1.0, 1))
	//fmt.Println(utils.COV(0.1, 1.0, 1))

	//na := utils.NewArray()
	//aa := "////HfrF7n8LhBb2KKNqrEpFEo//rubix/bacnet_server/points/ao/analogOutput-1/object_name"
	//s := strings.Split("analogOutput-1", "-")
	//fmt.Println(s)
	//
	//
	//mArr := utils.NewArray()
	//ss := strings.Split("analogOutput-1", "-")
	//for i, e := range ss {
	//		fmt.Println(e, i)
	//		if e != ""{
	//			mArr.Add(e)
	//		}
	//}
	//fmt.Println(mArr.Get(0))

	//
	//	if i == 10 { //AO
	//		if strings.Contains(e, "analogOutput") {
	//			fmt.Println("Yes")
	//			fmt.Println("analogOutput")
	//			fmt.Println("analogOutput", e)
	//			res := strings.ReplaceAll(e, "/", "")
	//			fmt.Println("analogOutput", res)
	//			//ss := strings.SplitAfter(res, "-")
	//			//fmt.Println("analogOutput", ss, 11111, res)
	//		}
	//		//fmt.Println(e, i)
	//	}
	//
	//
	//	//fmt.Println(e, i)
	//}



	//fmt.Println(aaa.Data()[2], 999)



	//fmt.Println(11111)
	//s := strings.SplitAfter(aa, "/")
	//na.Add(s)
	//fmt.Println(na.Values())
	//ss := strings.Split(aa, "/")
	//v := strings.Fields(aa)

}

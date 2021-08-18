package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)


func s(t string) string{
	switch t {
	case model.DriverTypeEnum.Serial:
		return model.DriverTypeEnum.Serial
	}
	return "nil"
}

func main(){
	fmt.Println(s("serial"))
}

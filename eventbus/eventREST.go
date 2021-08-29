package eventbus

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)

func EventREST(uuid string, body interface{}, ip string, port string, token string, write bool, thingType string) (interface{}, error) {
	c := client.NewSession("admin", "admin", "0.0.0.0", "1660")
	if thingType == model.CommonNaming.Point && write {
		point, err := c.ClientEditPoint(uuid, body)
		if err != nil {
			return nil, err
		}
		fmt.Println(point.Points.UUID)
		return point, err
	} else if thingType == model.CommonNaming.Point{
		point, err := c.ClientGetPoint(uuid)
		if err != nil {
			return nil, err
		}
		fmt.Println(point.Points.UUID)
		return point, err
	}
	return nil, nil

}



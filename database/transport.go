package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)

// GetTransports returns all items.
func (d *GormDatabase) GetTransports() ([]*interface{}, error) {
	var m []*interface{}
	query := d.DB.Find(&m);if query.Error != nil {
		return nil, query.Error
	}
	return m , nil
}


////ValidateTransport check the type, as in type=serial
//func validateTransport(body *model.TransportBody) (interface{}, error) {
//	t := body.TransportType
//	if t == model.CommonNaming.Serial {
//		var bk *model.SerialConnection
//		return bk, nil
//	}
//	return nil, nil
//
//}


// CreateTransport creates a thing.
func (d *GormDatabase) CreateTransport(t string, body interface{}) (interface{}, error) {
	fmt.Println(t, 9999, body)
	if t == model.CommonNaming.Serial {
		if x, ok := body.(model.SerialConnection); ok {
			fmt.Println(x.SerialPort, 9999999999)
		} else {
			fmt.Println("err", ok)
		}

	}
	return nil, nil
}

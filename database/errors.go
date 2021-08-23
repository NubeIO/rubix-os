package database

import (
	"errors"
	"fmt"
	"log"
)


func  errorMsg (field string, msg string, err error) error {
	r := fmt.Sprintf("error: %s - %s -%v", field, msg, err)
	log.Println(r)
	return errors.New(r)

}

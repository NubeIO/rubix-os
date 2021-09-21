package database

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func errorMsg(field string, msg string, err error) error {
	r := fmt.Sprintf("Error: %s - %s - %v", field, msg, err)
	log.Error(r)
	return errors.New(r)
}

package database

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func newDetailedError(field string, msg string, err error) error {
	r := fmt.Sprintf("Error: %s - %s - %v", field, msg, err)
	log.Error(r)
	return errors.New(r)
}

func newError(field string, msg string) error {
	r := fmt.Sprintf("Error: %s - %s", field, msg)
	log.Error(r)
	return errors.New(r)
}

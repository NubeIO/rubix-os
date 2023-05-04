package main

import (
	"errors"
	"fmt"
	"time"
)

func (inst *Instance) parseReadLoop(id string) *payloadReadPV {
	d, ok := inst.store.Get(id)
	if ok {
		parse := d.(*payloadReadPV)
		fmt.Println(1111, parse)
		return parse
	} else {
		return nil
	}
}

func (inst *Instance) readLoop(id string) (*payloadReadPV, error) {
	timeout := time.After(5 * time.Second)
	ticker := time.Tick(30 * time.Millisecond)
	for { // Keep trying until we're timed out or get a result/error
		select {
		case <-timeout: // Got a timeout! fail with a timeout error
			// maybe, check for one last time
			payload := inst.parseReadLoop(id)
			if payload == nil {
				return nil, errors.New("timed out")
			}
			return payload, nil
		case <-ticker:
			payload := inst.parseReadLoop(id)
			if payload == nil {
				return nil, nil
			} else if payload != nil {
				return payload, nil
			}
		}
	}
}
